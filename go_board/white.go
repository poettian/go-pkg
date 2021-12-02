package go_board

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

func StringToBytes(s []byte) {
	err := binary.Write(os.Stdout, binary.BigEndian, s)
	if err != nil {
		log.Panic(err)
	}
}

func CallMeTimer(d int) {
	timer := time.NewTimer(time.Second * 5)
	t := <-timer.C
	fmt.Println(t.Format("2006-01-02 03:04:05"))
}

func BufioScanner()  {
	input := "abcdefghijkl"
	scanner := bufio.NewScanner(strings.NewReader(input))
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) > 0 {
			return len(data), data, nil
		}
		if len(data) >= 5 {
			return 5, data[:5], nil
		}
		return 0, nil, nil
	}
	scanner.Split(split)
	buf := make([]byte, 2)
	scanner.Buffer(buf, bufio.MaxScanTokenSize)
	for scanner.Scan() {
		fmt.Printf("%s\n", scanner.Text())
	}
}

func ConfigFileExists() bool {
	configFile := "./config/app_dev.yaml"

	absPath,_ := filepath.Abs(configFile)
	fmt.Println(absPath)

	if _,err := os.Stat(configFile);os.IsNotExist(err) {
		return false
	}
	return true
}

func HandleSignal()  {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	time.Sleep(time.Second * 5)
	log.Println("goroutine signal handle")
}

func GetRandomToken(length int) string {
	r := make([]byte, length)
	io.ReadFull(rand.Reader, r)
	return base64.URLEncoding.EncodeToString(r)
}

func TimerTick()  {
	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case t := <-ticker.C:
				fmt.Println(t.Format("2006-01-02 15:04:05"))
			default:
				fmt.Println("default")
			}
			time.Sleep(time.Millisecond * 200)
		}
	}()
	<-time.After(time.Second * 10)
	fmt.Println("game over")
}

func GetCpuNums() {
	n := runtime.NumCPU()
	fmt.Printf("cpu num is: %d\n", n)
}

func ListenAndServe(address string) {
	lister, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("listen err: %v", err)
	}
	defer lister.Close()
	log.Printf("bind:%s, start listening...", address)

	for {
		// Accept 会一直阻塞直到有新的连接建立或者listen中断才会返回
		conn, err := lister.Accept()
		if err != nil {
			// 通常是由于listener被关闭无法继续监听导致的错误
			log.Fatalf("accept err: %v", err)
		}
		go HandleConn(conn)
	}
}

func HandleConn(conn net.Conn) {
	// 使用 bufio 标准库提供的缓冲区功能
	reader := bufio.NewReader(conn)
	for {
		// ReadString 会一直阻塞直到遇到分隔符 '\n'
		// 遇到分隔符后会返回上次遇到分隔符或连接建立后收到的所有数据, 包括分隔符本身
		// 若在遇到分隔符之前遇到异常, ReadString 会返回已收到的数据和错误信息
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("connection close")
			} else {
				log.Println(err)
			}
			return
		}
		b := []byte(msg)
		// 将收到的信息发送给客户端
		conn.Write(b)
	}
}