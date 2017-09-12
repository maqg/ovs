package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"octlink/ovs/utils/config"
	"octlink/ovs/utils/octlog"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/pkg/errors"
)

const (
	timeStrFormat = "2006-01-02 15:04:05"
)

// Time2Str convert int64 time to string
func Time2Str(timeVal int64) string {
	return time.Unix(timeVal, 0).Format(timeStrFormat)
}

// Time2StrSimple convert int64 time to string
func Time2StrSimple(timeVal int64) string {
	return time.Unix(timeVal, 0).Format("20060102150405")
}

// CurrentTime get current time in int64 format
func CurrentTime() int64 {
	return int64(time.Now().Unix())
}

// CurrentTimeStr get current time in string format
func CurrentTimeStr() string {
	return Time2Str(CurrentTime())
}

// CurrentTimeSimple get current time in string format
func CurrentTimeSimple() string {
	return Time2StrSimple(CurrentTime())
}

// Version of this program
func Version() string {
	return "0.0.1"
}

// IsFileExist check file's existence
func IsFileExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// ParseInt to parse int
func ParseInt(val interface{}) int {
	return int(val.(float64))
}

// ParseInt64 to parse int
func ParseInt64(val interface{}) int64 {
	return int64(val.(float64))
}

// StringToBytes return GoString's buffer slice(enable modify string)
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

// BytesToString convert b to string without copy
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StringToPointer returns &s[0], which is not allowed in go
func StringToPointer(s string) unsafe.Pointer {
	p := (*reflect.StringHeader)(unsafe.Pointer(&s))
	return unsafe.Pointer(p.Data)
}

// BytesToPointer returns &b[0], which is not allowed in go
func BytesToPointer(b []byte) unsafe.Pointer {
	p := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	return unsafe.Pointer(p.Data)
}

// GetDigest return digest hash value of sha256
func GetDigest(buffer []byte) string {
	hash := sha256.New()
	hash.Write(buffer)
	return hex.EncodeToString(hash.Sum(nil))
}

// GetDigestStr get digest value by string
func GetDigestStr(data string) string {
	return GetDigest(StringToBytes(data))
}

// CreateDir create dir if not exist
func CreateDir(filepath string) {
	if !IsFileExist(filepath) {
		os.MkdirAll(filepath, os.ModePerm)
	}
}

// RemoveDir if file or directory exists, just remove it
func RemoveDir(filepath string) {
	if IsFileExist(filepath) {
		os.RemoveAll(filepath)
	}
}

// RemoveFile if file or directory exists, just remove it
func RemoveFile(filepath string) {
	if IsFileExist(filepath) {
		os.Remove(filepath)
	}
}

// JSON2String convert json to string
func JSON2String(v interface{}) string {
	data, _ := json.MarshalIndent(v, "", "  ")
	return string(data)
}

// JSON2Bytes convert json to string
func JSON2Bytes(v interface{}) []byte {
	data, _ := json.MarshalIndent(v, "", "  ")
	return data
}

// TrimDir trim directory name by replace "//"
func TrimDir(path string) string {
	return strings.Replace(path, "//", "/", -1)
}

// GetFileSize to get file length
func GetFileSize(path string) int64 {
	stat, err := os.Stat(path)
	if err != nil {
		return 0
	}

	return stat.Size()
}

// OCTSystem for syscal command calling
func OCTSystem(cmdstr string) (string, error) {

	cmd := exec.Command("/bin/sh", "-c", cmdstr)
	data, err := cmd.Output()
	if err != nil {
		fmt.Printf("get cmd output error of %s,%s\n", cmdstr, err)
		return "", err
	}

	return BytesToString(data), nil
}

// StringToInt convert string to int value
func StringToInt(src string) int {
	ret, err := strconv.Atoi(src)
	if err != nil {
		return -1
	}
	return ret
}

// StringToInt64 convert string to int value
func StringToInt64(src string) int64 {
	src = strings.Replace(src, "\n", "", -1)
	ret, err := strconv.ParseInt(src, 10, 64)
	if err != nil {
		return -1
	}
	return ret
}

// IntToString convert int to string value
func IntToString(src int) string {
	return strconv.Itoa(src)
}

// Int64ToString convert int64 to string value
func Int64ToString(src int64) string {
	return strconv.FormatInt(src, 10)
}

// FileToBytes for filepath convert to bytes
func FileToBytes(filepath string) []byte {
	if !IsFileExist(filepath) {
		return nil
	}

	fd, err := os.Open(filepath)
	if err != nil {
		return nil
	}

	defer fd.Close()

	data, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil
	}

	return data
}

// FileToString convert file content to string
func FileToString(filepath string) string {
	return BytesToString(FileToBytes(filepath))
}

// NumberToInt convert int,int32,int64,float,float32, to int
func NumberToInt(value interface{}) int {
	switch reflect.TypeOf(value).Kind() {
	case reflect.Int:
		return value.(int)
	case reflect.Int64:
		return int(value.(int64))
	case reflect.Float32:
		return int(value.(float32))
	case reflect.Float64:
		return int(value.(float64))
	}
	return 0
}

// CopyFile for srcfile to dst file, return size on success
func CopyFile(srcFile, dstFile string) (int64, error) {
	sd, err := os.Open(srcFile)
	if err != nil {
		octlog.Error("open src file of %s error: %s\n", srcFile, err)
		return 0, err
	}

	defer sd.Close()

	dd, err := os.Create(dstFile)
	if err != nil {
		octlog.Error("open dst file of %s error: %s\n", dstFile, err)
		return 0, err
	}
	defer dd.Close()

	return io.Copy(dd, sd)
}

// OSType return os type
func OSType() string {
	// can be darwin,windows,linux
	return runtime.GOOS
}

// IsPlatformWindows for platform type judgement
func IsPlatformWindows() bool {
	return OSType() == config.OSTypeWindows
}

var logger *octlog.LogConfig

// InitLog to init api log config
func InitLog(level int) {
	logger = octlog.InitLogConfig("utils.log", level)
}

// RemoveFromArray for remove from array
func RemoveFromArray(slice []interface{}, s int) []interface{} {
	return append(slice[:s], slice[s+1:]...)
}

// RemoveFromArrayEx for remove from array without seq indexing
func RemoveFromArrayEx(s []interface{}, i int) []interface{} {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

// JSONDecodeHTTPRequest for json decode http request
func JSONDecodeHTTPRequest(req *http.Request, val interface{}) (err error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return errors.Wrap(err, "unable to read the request, %s")
	}

	if err = json.Unmarshal(body, val); err != nil {
		return errors.Wrap(err, fmt.Sprintf("unable to parse string '%s' to JSON object", string(body)))
	}

	return nil
}

// InetAton Convert net.IP to int64
func InetAton(addr string) int64 {
	bits := strings.Split(addr, ".")
	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])
	return int64(b0)<<24 + int64(b1)<<16 + int64(b2)<<8 + int64(b3)
}

// InetNtoa convert int64 address to string like "xx.xx.xx.xx"
func InetNtoa(ipnr int64) string {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0]).String()
}

// CIDRToNetmask convert 16 to "255.255.0.0"
func CIDRToNetmask(cidr int) string {

	var ipint int64

	if cidr > 32 {
		cidr = 32
	}

	for i := 0; i < cidr; i++ {
		ipint += 1 << uint(32-i-1)
	}

	return InetNtoa(ipint)
}

// NetmaskToCIDR change "255.255.0.0" to 16
func NetmaskToCIDR(netmask string) int {
	countBit := func(num uint) int {
		count := uint(0)
		var i uint
		for i = 31; i > 0; i-- {
			count += ((num << i) >> uint(31)) & uint(1)
		}

		return int(count)
	}

	cidr := 0
	for _, o := range strings.Split(netmask, ".") {
		num, err := strconv.ParseUint(o, 10, 32)
		if err != nil {
			return -1
		}
		cidr += countBit(uint(num))
	}

	return cidr
}

// GetNetworkNumber get network number by ip and mask
func GetNetworkNumber(ip, netmask string) (string, error) {
	ips := strings.Split(ip, ".")
	masks := strings.Split(netmask, ".")

	ipInByte := make([]interface{}, 4)
	for i := 0; i < len(ips); i++ {
		p, err := strconv.ParseUint(ips[i], 10, 32)
		if err != nil {
			logger.Panicf("unable to get network number[ip:%v, netmask:%v], error %s\n", ip, netmask, err)
			return "", errors.Wrap(err, fmt.Sprintf("unable to get network number[ip:%v, netmask:%v]", ip, netmask))
		}
		m, err := strconv.ParseUint(masks[i], 10, 32)
		if err != nil {
			logger.Panicf("unable to get network number[ip:%v, netmask:%v], error %s\n", ip, netmask, err)
			return "", errors.Wrap(err, fmt.Sprintf("unable to get network number[ip:%v, netmask:%v]", ip, netmask))
		}
		ipInByte[i] = p & m
	}

	cidr := NetmaskToCIDR(netmask)
	if cidr == -1 {
		return "", fmt.Errorf("unable to get network number[ip:%v, netmask:%v]", ip, netmask)
	}

	return fmt.Sprintf("%v.%v.%v.%v/%v", ipInByte[0], ipInByte[1], ipInByte[2], ipInByte[3], cidr), nil
}
