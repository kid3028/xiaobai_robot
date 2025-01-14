package solia

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/axgle/mahonia"
	"io"
	"math/rand"
	"os"
	"strings"
	"unicode/utf8"
)

// 空结构体
var Exists = struct{}{}

//Set
type Set struct {
	// struct为结构体类型的变量
	m map[interface{}]struct{}
}

func (s *Set) Add(items ...interface{}) error {
	for _, item := range items {
		s.m[item] = Exists
	}
	return nil
}

func (s *Set) Contains(item interface{}) bool {
	_, ok := s.m[item]
	return ok
}

// 清空set
func (s *Set) Clear() {
	s.m = make(map[interface{}]struct{})
}

type Solia struct {
	UserId string //正在接龙的用户id
	StrSet Set    //已经使用过的成语
	flag   bool   //是否开始接龙
	rd     *bufio.Reader
	tryNum int    // 尝试次数
	nowStr string //当前接龙的成语
}

func (s *Solia) ReNew() {
	s.tryNum = 3
	s.StrSet.Clear()
	s.UserId = ""
	s.nowStr = ""
}

func (s *Solia) ReadStart(userID string) (string, error) {
	//打开文件
	s.tryNum = 3
	s.StrSet.m = make(map[interface{}]struct{})
	if s.UserId == "" {
		s.UserId = userID
	} else {
		if s.UserId != userID {
			return "", errors.New(Wait)
		} else {
			return "", errors.New(ReStart)
		}
	}
	n := rand.Intn(9999) + 1
	str, err := s.readLineNum(n)
	s.StrSet.Add(str)
	s.nowStr = str
	return str, err
}

func (s *Solia) readLineNum(lineNum int) (string, error) {
	num := 1
	strPath, _ := os.Getwd()
	fmt.Println(strPath[:strings.LastIndex(strPath, "robot")+5])
	strPath = strPath[:strings.LastIndex(strPath, "robot")+5]
	file, err := os.Open(strPath + "./idiom.txt") //只是用来读的时候，用os.Open。绝对路径，获取robot路径
	if err != nil {
		fmt.Printf("打开文件失败,err:%v\n", err)
		return "", err
	}
	defer file.Close()                   //关闭文件,为了避免文件泄露和忘记写关闭文件
	decoder := mahonia.NewDecoder("gbk") //转码，避免中文字符乱码
	//使用buffio读取文件内容
	reader := bufio.NewReader(decoder.NewReader(file)) //创建新的读的对象
	for {
		line, err := reader.ReadString('\n') //注意是字符，换行符。
		if err == io.EOF {
			fmt.Println("文件读完了")
			break
		}
		if err != nil { //错误处理
			fmt.Printf("读取文件失败,错误为:%v", err)
			return "", err
		}
		num++
		line = strings.Replace(line, "\r\n", "", -1)
		if num >= lineNum && utf8.RuneCountInString(line) == 4 {
			return line, nil
		}
		//fmt.Println(utf8.RuneCountInString(line),string(line))
	}
	return "人山人海", nil
}

func (s *Solia) ReadStr(content string) (string, error) {
	if s.StrSet.Contains(content) { //判断是否重复的成语
		if s.tryNum > 1 {
			s.tryNum--
			return "", errors.New(fmt.Sprintf(ContainsNotOver, s.tryNum))
		} else {
			s.ReNew()
			return "", errors.New(fmt.Sprintf(ContainsOver3))
		}
	}
	strs := strings.Split(content, "> ")
	content = strs[len(strs)-1]
	s1 := string([]rune(s.nowStr)[3:])
	if s1 != string([]rune(content)[:1]) { //判断首字是否和上一个成语尾字一样
		if s.tryNum > 1 {
			s.tryNum--
			return "", errors.New(fmt.Sprintf(SameNotOver, s.tryNum))
		} else {
			s.ReNew()
			return "", errors.New(fmt.Sprintf(SameNotOver3))
		}
	}
	var res string
	var flag bool
	if len([]rune(content)) >= 4 {
		str1 := string([]rune(content)[3:])

		strPath, _ := os.Getwd()
		fmt.Println(strPath[:strings.LastIndex(strPath, "robot")+5])
		strPath = strPath[:strings.LastIndex(strPath, "robot")+5]
		file, err := os.Open(strPath + "./idiom.txt") //只是用来读的时候，用os.Open。绝对路径，获取robot路径
		if err != nil {
			fmt.Printf("打开文件失败,err:%v\n", err)
			return "", err
		}
		defer file.Close()                   //关闭文件,为了避免文件泄露和忘记写关闭文件
		decoder := mahonia.NewDecoder("gbk") //转码，避免中文字符乱码
		//使用buffio读取文件内容
		reader := bufio.NewReader(decoder.NewReader(file)) //创建新的读的对象
		for {
			line, err := reader.ReadString('\n') //注意是字符，换行符。
			if err == io.EOF {
				fmt.Println("文件读完了")
				break
			}
			if err != nil { //错误处理
				fmt.Printf("读取文件失败,错误为:%v", err)
				s.ReNew()
				return "", errors.New("我不会了，你赢了")
			}
			line = strings.Replace(line, "\r\n", "", -1)
			if str1 == string([]rune(line)[:1]) && !s.StrSet.Contains(line) {
				res = line
			}
			if line == content {
				flag = true
			}
			if flag && res != "" {
				break
			}
			//fmt.Println(utf8.RuneCountInString(line),string(line))
		}
	}
	if !flag {
		if s.tryNum > 1 {
			s.tryNum--
			return "", errors.New(fmt.Sprintf(ErrorNotOver, s.tryNum))
		} else {
			s.ReNew()
			return "", errors.New(fmt.Sprintf(ErrorOver3))
		}
	} else {
		//s.UserId=""
		//s.StrSet.Clear()
		s.StrSet.Add(content)
		s.nowStr = content
	}
	if res == "" {
		s.ReNew()
		return "", errors.New(fmt.Sprintf(Success))
	} else {
		s.StrSet.Add(res)
		s.nowStr = res
	}
	return res, nil
}
