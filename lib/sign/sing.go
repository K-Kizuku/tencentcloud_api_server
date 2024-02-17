package sign

import (
	"bytes"
	"compress/zlib"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

func GenUserSig(sdkappid int, key string, userid string, expire int) (string, error) {
	return genSig(sdkappid, key, userid, expire, nil)
}

func GenUserSigWithBuf(sdkappid int, key string, userid string, expire int, buf []byte) (string, error) {
	return genSig(sdkappid, key, userid, expire, buf)
}

func GenPrivateMapKey(sdkappid int, key string, userid string, expire int, roomid uint32, privilegeMap uint32) (string, error) {
	var userbuf []byte = genUserBuf(userid, sdkappid, roomid, expire, privilegeMap, 0, "")
	return genSig(sdkappid, key, userid, expire, userbuf)
}

func GenPrivateMapKeyWithStringRoomID(sdkappid int, key string, userid string, expire int, roomStr string, privilegeMap uint32) (string, error) {
	var userbuf []byte = genUserBuf(userid, sdkappid, 0, expire, privilegeMap, 0, roomStr)
	return genSig(sdkappid, key, userid, expire, userbuf)
}

func genUserBuf(account string, dwSdkappid int, dwAuthID uint32,
	dwExpTime int, dwPrivilegeMap uint32, dwAccountType uint32, roomStr string) []byte {

	offset := 0
	length := 1 + 2 + len(account) + 20 + len(roomStr)
	if len(roomStr) > 0 {
		length = length + 2
	}

	userBuf := make([]byte, length)

	//ver
	if len(roomStr) > 0 {
		userBuf[offset] = 1
	} else {
		userBuf[offset] = 0
	}

	offset++
	userBuf[offset] = (byte)((len(account) & 0xFF00) >> 8)
	offset++
	userBuf[offset] = (byte)(len(account) & 0x00FF)
	offset++

	for ; offset < len(account)+3; offset++ {
		userBuf[offset] = account[offset-3]
	}

	//dwSdkAppid
	userBuf[offset] = (byte)((dwSdkappid & 0xFF000000) >> 24)
	offset++
	userBuf[offset] = (byte)((dwSdkappid & 0x00FF0000) >> 16)
	offset++
	userBuf[offset] = (byte)((dwSdkappid & 0x0000FF00) >> 8)
	offset++
	userBuf[offset] = (byte)(dwSdkappid & 0x000000FF)
	offset++

	//dwAuthId
	userBuf[offset] = (byte)((dwAuthID & 0xFF000000) >> 24)
	offset++
	userBuf[offset] = (byte)((dwAuthID & 0x00FF0000) >> 16)
	offset++
	userBuf[offset] = (byte)((dwAuthID & 0x0000FF00) >> 8)
	offset++
	userBuf[offset] = (byte)(dwAuthID & 0x000000FF)
	offset++

	//dwExpTime now+300;
	currTime := time.Now().Unix()
	var expire = currTime + int64(dwExpTime)
	userBuf[offset] = (byte)((expire & 0xFF000000) >> 24)
	offset++
	userBuf[offset] = (byte)((expire & 0x00FF0000) >> 16)
	offset++
	userBuf[offset] = (byte)((expire & 0x0000FF00) >> 8)
	offset++
	userBuf[offset] = (byte)(expire & 0x000000FF)
	offset++

	//dwPrivilegeMap
	userBuf[offset] = (byte)((dwPrivilegeMap & 0xFF000000) >> 24)
	offset++
	userBuf[offset] = (byte)((dwPrivilegeMap & 0x00FF0000) >> 16)
	offset++
	userBuf[offset] = (byte)((dwPrivilegeMap & 0x0000FF00) >> 8)
	offset++
	userBuf[offset] = (byte)(dwPrivilegeMap & 0x000000FF)
	offset++

	//dwAccountType
	userBuf[offset] = (byte)((dwAccountType & 0xFF000000) >> 24)
	offset++
	userBuf[offset] = (byte)((dwAccountType & 0x00FF0000) >> 16)
	offset++
	userBuf[offset] = (byte)((dwAccountType & 0x0000FF00) >> 8)
	offset++
	userBuf[offset] = (byte)(dwAccountType & 0x000000FF)
	offset++

	if len(roomStr) > 0 {
		userBuf[offset] = (byte)((len(roomStr) & 0xFF00) >> 8)
		offset++
		userBuf[offset] = (byte)(len(roomStr) & 0x00FF)
		offset++

		for ; offset < length; offset++ {
			userBuf[offset] = roomStr[offset-(length-len(roomStr))]
		}
	}

	return userBuf
}

func hmacsha256(sdkappid int, key string, identifier string, currTime int64, expire int, base64UserBuf *string) string {
	var contentToBeSigned string
	contentToBeSigned = "TLS.identifier:" + identifier + "\n"
	contentToBeSigned += "TLS.sdkappid:" + strconv.Itoa(sdkappid) + "\n"
	contentToBeSigned += "TLS.time:" + strconv.FormatInt(currTime, 10) + "\n"
	contentToBeSigned += "TLS.expire:" + strconv.Itoa(expire) + "\n"
	if nil != base64UserBuf {
		contentToBeSigned += "TLS.userbuf:" + *base64UserBuf + "\n"
	}

	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(contentToBeSigned))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func genSig(sdkappid int, key string, identifier string, expire int, userbuf []byte) (string, error) {
	currTime := time.Now().Unix()
	sigDoc := make(map[string]interface{})
	sigDoc["TLS.ver"] = "2.0"
	sigDoc["TLS.identifier"] = identifier
	sigDoc["TLS.sdkappid"] = sdkappid
	sigDoc["TLS.expire"] = expire
	sigDoc["TLS.time"] = currTime
	var base64UserBuf string
	if nil != userbuf {
		base64UserBuf = base64.StdEncoding.EncodeToString(userbuf)
		sigDoc["TLS.userbuf"] = base64UserBuf
		sigDoc["TLS.sig"] = hmacsha256(sdkappid, key, identifier, currTime, expire, &base64UserBuf)
	} else {
		sigDoc["TLS.sig"] = hmacsha256(sdkappid, key, identifier, currTime, expire, nil)
	}

	data, err := json.Marshal(sigDoc)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	if _, err = w.Write(data); err != nil {
		return "", err
	}
	if err = w.Close(); err != nil {
		return "", err
	}
	return base64urlEncode(b.Bytes()), nil
}

func VerifyUserSig(sdkappid uint64, key string, userid string, usersig string, now time.Time) error {
	sig, err := newUserSig(usersig)
	if err != nil {
		return err
	}
	return sig.verify(sdkappid, key, userid, now, nil)
}

func VerifyUserSigWithBuf(sdkappid uint64, key string, userid string, usersig string, now time.Time, userbuf []byte) error {
	sig, err := newUserSig(usersig)
	if err != nil {
		return err
	}
	return sig.verify(sdkappid, key, userid, now, userbuf)
}

type userSig struct {
	Version    string `json:"TLS.ver,omitempty"`
	Identifier string `json:"TLS.identifier,omitempty"`
	SdkAppID   uint64 `json:"TLS.sdkappid,omitempty"`
	Expire     int64  `json:"TLS.expire,omitempty"`
	Time       int64  `json:"TLS.time,omitempty"`
	UserBuf    []byte `json:"TLS.userbuf,omitempty"`
	Sig        string `json:"TLS.sig,omitempty"`
}

func newUserSig(usersig string) (userSig, error) {
	b, err := base64urlDecode(usersig)
	if err != nil {
		return userSig{}, err
	}
	r, err := zlib.NewReader(bytes.NewReader(b))
	if err != nil {
		return userSig{}, err
	}
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return userSig{}, err
	}
	if err = r.Close(); err != nil {
		return userSig{}, err
	}
	var sig userSig
	if err = json.Unmarshal(data, &sig); err != nil {
		return userSig{}, nil
	}
	return sig, nil
}

func (u userSig) verify(sdkappid uint64, key string, userid string, now time.Time, userbuf []byte) error {
	if sdkappid != u.SdkAppID {
		return ErrSdkAppIDNotMatch
	}
	if userid != u.Identifier {
		return ErrIdentifierNotMatch
	}
	if now.Unix() > u.Time+u.Expire {
		return ErrExpired
	}
	if userbuf != nil {
		if u.UserBuf == nil {
			return ErrUserBufTypeNotMatch
		}
		if !bytes.Equal(userbuf, u.UserBuf) {
			return ErrUserBufNotMatch
		}
	} else if u.UserBuf != nil {
		return ErrUserBufTypeNotMatch
	}
	if u.sign(key) != u.Sig {
		return ErrSigNotMatch
	}
	return nil
}

func (u userSig) sign(key string) string {
	var sb bytes.Buffer
	sb.WriteString("TLS.identifier:")
	sb.WriteString(u.Identifier)
	sb.WriteString("\n")
	sb.WriteString("TLS.sdkappid:")
	sb.WriteString(strconv.FormatUint(u.SdkAppID, 10))
	sb.WriteString("\n")
	sb.WriteString("TLS.time:")
	sb.WriteString(strconv.FormatInt(u.Time, 10))
	sb.WriteString("\n")
	sb.WriteString("TLS.expire:")
	sb.WriteString(strconv.FormatInt(u.Expire, 10))
	sb.WriteString("\n")
	if u.UserBuf != nil {
		sb.WriteString("TLS.userbuf:")
		sb.WriteString(base64.StdEncoding.EncodeToString(u.UserBuf))
		sb.WriteString("\n")
	}

	h := hmac.New(sha256.New, []byte(key))
	h.Write(sb.Bytes())
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// 错误类型
var (
	ErrSdkAppIDNotMatch    = errors.New("sdk appid not match")
	ErrIdentifierNotMatch  = errors.New("identifier not match")
	ErrExpired             = errors.New("expired")
	ErrUserBufTypeNotMatch = errors.New("userbuf type not match")
	ErrUserBufNotMatch     = errors.New("userbuf not match")
	ErrSigNotMatch         = errors.New("sig not match")
)

func base64urlEncode(data []byte) string {
	str := base64.StdEncoding.EncodeToString(data)
	str = strings.Replace(str, "+", "*", -1)
	str = strings.Replace(str, "/", "-", -1)
	str = strings.Replace(str, "=", "_", -1)
	return str
}

func base64urlDecode(str string) ([]byte, error) {
	str = strings.Replace(str, "_", "=", -1)
	str = strings.Replace(str, "-", "/", -1)
	str = strings.Replace(str, "*", "+", -1)
	return base64.StdEncoding.DecodeString(str)
}
