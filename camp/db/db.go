package db

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mufe/golang-base/camp/xlog"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type ViaSSHDialer struct {
	client *ssh.Client
	_      *context.Context
}

func (self *ViaSSHDialer) Dial(context context.Context, addr string) (net.Conn, error) {
	return self.client.Dial("tcp", addr)
}

type remoteScriptType byte
type remoteShellType byte

const (
	cmdLine remoteScriptType = iota
	rawScript
	scriptFile

	interactiveShell remoteShellType = iota
	nonInteractiveShell
)

type Db struct {
	db    *sql.DB
	Print bool
}

type Client struct {
	client *ssh.Client
}

// DialWithPasswd starts a client connection to the given SSH server with passwd authmethod.
func DialWithPasswd(addr, user, passwd string) (*Client, error) {

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(passwd),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}

	return Dial("tcp", addr, config)
}

// DialWithKey starts a client connection to the given SSH server with key authmethod.
func DialWithKey(addr, user, keyfile string) (*Client, error) {
	key, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}

	return Dial("tcp", addr, config)
}

// DialWithKeyWithPassphrase same as DialWithKey but with a passphrase to decrypt the private key
func DialWithKeyWithPassphrase(addr, user, keyfile string, passphrase string) (*Client, error) {
	key, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKeyWithPassphrase(key, []byte(passphrase))
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}

	return Dial("tcp", addr, config)
}

// Dial starts a client connection to the given SSH server.
// This is wrap the ssh.Dial
func Dial(network, addr string, config *ssh.ClientConfig) (*Client, error) {
	client, err := ssh.Dial(network, addr, config)
	if err != nil {
		return nil, err
	}
	return &Client{
		client: client,
	}, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}

// Cmd create a command on client
func (c *Client) Cmd(cmd string) *remoteScript {
	return &remoteScript{
		_type:  cmdLine,
		client: c.client,
		script: bytes.NewBufferString(cmd + "\n"),
	}
}

// Script
func (c *Client) Script(script string) *remoteScript {
	return &remoteScript{
		_type:  rawScript,
		client: c.client,
		script: bytes.NewBufferString(script + "\n"),
	}
}

// ScriptFile
func (c *Client) ScriptFile(fname string) *remoteScript {
	return &remoteScript{
		_type:      scriptFile,
		client:     c.client,
		scriptFile: fname,
	}
}

type remoteScript struct {
	client     *ssh.Client
	_type      remoteScriptType
	script     *bytes.Buffer
	scriptFile string
	err        error

	stdout io.Writer
	stderr io.Writer
}

// Run
func (rs *remoteScript) Run() error {
	if rs.err != nil {
		fmt.Println(rs.err)
		return rs.err
	}

	if rs._type == cmdLine {
		return rs.runCmds()
	} else if rs._type == rawScript {
		return rs.runScript()
	} else if rs._type == scriptFile {
		return rs.runScriptFile()
	} else {
		return errors.New("Not supported remoteScript type")
	}
}

func (rs *remoteScript) Output() ([]byte, error) {
	if rs.stdout != nil {
		return nil, errors.New("Stdout already set")
	}
	var out bytes.Buffer
	rs.stdout = &out
	err := rs.Run()
	return out.Bytes(), err
}

func (rs *remoteScript) SmartOutput() ([]byte, error) {
	if rs.stdout != nil {
		return nil, errors.New("Stdout already set")
	}
	if rs.stderr != nil {
		return nil, errors.New("Stderr already set")
	}

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)
	rs.stdout = &stdout
	rs.stderr = &stderr
	err := rs.Run()
	if err != nil {
		return stderr.Bytes(), err
	}
	return stdout.Bytes(), err
}

func (rs *remoteScript) Cmd(cmd string) *remoteScript {
	_, err := rs.script.WriteString(cmd + "\n")
	if err != nil {
		rs.err = err
	}
	return rs
}

func (rs *remoteScript) SetStdio(stdout, stderr io.Writer) *remoteScript {
	rs.stdout = stdout
	rs.stderr = stderr
	return rs
}

func (rs *remoteScript) runCmd(cmd string) error {
	session, err := rs.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	session.Stdout = rs.stdout
	session.Stderr = rs.stderr

	if err := session.Run(cmd); err != nil {
		return err
	}
	return nil
}

func (rs *remoteScript) runCmds() error {
	for {
		statment, err := rs.script.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if err := rs.runCmd(statment); err != nil {
			return err
		}
	}

	return nil
}

func (rs *remoteScript) runScript() error {
	session, err := rs.client.NewSession()
	if err != nil {
		return err
	}

	session.Stdin = rs.script
	session.Stdout = rs.stdout
	session.Stderr = rs.stderr

	if err := session.Shell(); err != nil {
		return err
	}
	if err := session.Wait(); err != nil {
		return err
	}

	return nil
}

func (rs *remoteScript) runScriptFile() error {
	var buffer bytes.Buffer
	file, err := os.Open(rs.scriptFile)
	if err != nil {
		return err
	}
	_, err = io.Copy(&buffer, file)
	if err != nil {
		return err
	}

	rs.script = &buffer
	return rs.runScript()
}

type remoteShell struct {
	client         *ssh.Client
	requestPty     bool
	terminalConfig *TerminalConfig

	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

type TerminalConfig struct {
	Term   string
	Hight  int
	Weight int
	Modes  ssh.TerminalModes
}

// Terminal create a interactive shell on client.
func (c *Client) Terminal(config *TerminalConfig) *remoteShell {
	return &remoteShell{
		client:         c.client,
		terminalConfig: config,
		requestPty:     true,
	}
}

// Shell create a noninteractive shell on client.
func (c *Client) Shell() *remoteShell {
	return &remoteShell{
		client:     c.client,
		requestPty: false,
	}
}

func (rs *remoteShell) SetStdio(stdin io.Reader, stdout, stderr io.Writer) *remoteShell {
	rs.stdin = stdin
	rs.stdout = stdout
	rs.stderr = stderr
	return rs
}

// Start start a remote shell on client
func (rs *remoteShell) Start() error {
	session, err := rs.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	if rs.stdin == nil {
		session.Stdin = os.Stdin
	} else {
		session.Stdin = rs.stdin
	}
	if rs.stdout == nil {
		session.Stdout = os.Stdout
	} else {
		session.Stdout = rs.stdout
	}
	if rs.stderr == nil {
		session.Stderr = os.Stderr
	} else {
		session.Stderr = rs.stderr
	}

	if rs.requestPty {
		tc := rs.terminalConfig
		if tc == nil {
			tc = &TerminalConfig{
				Term:   "xterm",
				Hight:  40,
				Weight: 80,
			}
		}
		if err := session.RequestPty(tc.Term, tc.Hight, tc.Weight, tc.Modes); err != nil {
			return err
		}
	}

	if err := session.Shell(); err != nil {
		return err
	}

	if err := session.Wait(); err != nil {
		return err
	}

	return nil
}

func (s *Db) Connect(host string, port int, user, pwd, database string) error {
	var err error
	path := strings.Join([]string{user, ":", pwd, "@tcp(", host, ":", strconv.Itoa(port), ")/", database, "?charset=utf8mb4"}, "")
	s.db, err = sql.Open("mysql", path)
	if err != nil {
		xlog.ErrorP(err)
		return err
	}
	if s.Print {
		err := mysql.SetLogger(log.New(os.Stdout, "[MySQL] ", log.LstdFlags|log.Lshortfile))
		if err != nil {
			return err
		}
	}
	s.db.SetMaxOpenConns(200)
	s.db.SetMaxIdleConns(100)
	return nil
}

func (s *Db) Query(sql string, args ...interface{}) (rows *sql.Rows, err error) {
	t := time.Now()
	rows, err = s.db.Query(sql, args...)
	if s.Print {
		xlog.DB(false, time.Now().Sub(t), 0, sql, args...)
	}
	if err != nil {
		xlog.ErrorP(err)
	}
	return rows, err
}

func (s *Db) QueryRow(sql string, args ...interface{}) (result *sql.Row) {
	t := time.Now()
	result = s.db.QueryRow(sql, args...)
	if s.Print {
		xlog.DB(false, time.Now().Sub(t), 0, sql, args...)
	}
	return result
}

func (s *Db) Exec(sql string, args ...interface{}) (result sql.Result, err error) {
	t := time.Now()
	result, err = s.db.Exec(sql, args...)
	var affected int64
	if err == nil {
		affected, _ = result.RowsAffected()
	} else {
		xlog.ErrorP(err)
	}
	if s.Print {
		xlog.DB(true, time.Now().Sub(t), affected, sql, args...)
	}
	return result, err
}

func (s *Db) WithTransaction(f func(tx *Tx) error) error {
	var t = &Tx{
		Print: s.Print,
	}
	var err error
	t.tx, err = s.db.Begin()
	if err != nil {
		xlog.ErrorP(err)
		return err
	}
	// 调用f时如果出现panic，err则会无法正常赋值，因此需要此变量
	var success bool
	defer func() {
		if !success {
			// 执行f时出现任何问题，都要Rollback
			t.tx.Rollback()
		}
	}()
	err = f(t)
	if err != nil {
		xlog.ErrorP(err)
		return err
	}
	// 提交
	err = t.tx.Commit()
	if err != nil {
		xlog.ErrorP(err)
		return err
	}
	success = true
	return nil
}
