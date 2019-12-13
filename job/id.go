package job // import "gopkg.in/webnice/job.v1/job"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// Parse Разбирается переданная строка в объект ID
func (ido *ID) Parse(src string) (ret *ID, err error) {
	var tmp []string

	ret = new(ID)
	tmp = strings.Split(src, ":")
	ret.name = tmp[0]
	if len(tmp) >= 2 {
		ret.serialNumber, _ = strconv.ParseUint(tmp[1], 10, 64)
	}
	if len(tmp) >= 3 {
		ret.pid, _ = strconv.ParseInt(tmp[2], 10, 64)
	}

	return
}

// String Возвращает строковый эквивалент ID процесса
func (ido *ID) String() string { return fmt.Sprintf("%s:%d:%d", ido.name, ido.serialNumber, ido.pid) }

// EqualRoot Равны по корню.
// Сравниваются: пакет, имя структуры
// Не учитываются: PID, порядковый номер
func (ido *ID) EqualRoot(id string) bool {
	if src, err := ido.Parse(id); err == nil &&
		strings.EqualFold(ido.name, src.name) {
		return true
	}
	return false
}

// EqualParent Равны в пределах приложения.
// Сравниваются: пакет, имя структуры, PID
// Не учитывается: порядковый номер
func (ido *ID) EqualParent(id string) bool {
	if src, err := ido.Parse(id); err == nil &&
		strings.EqualFold(ido.name, src.name) &&
		ido.pid == src.pid {
		return true
	}
	return false
}

// Equal Равны в пределах пакета
// Сравниваются: пакет, имя структуры, порядковый номер
// Не учитывается: PID
func (ido *ID) Equal(id string) bool {
	if src, err := ido.Parse(id); err == nil &&
		strings.EqualFold(ido.name, src.name) &&
		ido.serialNumber == src.serialNumber {
		return true
	}
	return false
}

// EqualFold Полное сравнение.
// Сравниваются: пакет, имя структуры, порядковый номер, PID
func (ido *ID) EqualFold(id string) bool { return strings.EqualFold(ido.String(), id) }

// Получение уникального имени пакета + имя структуры
func getStructName(obj interface{}) string {
	var (
		rt                         reflect.Type
		packageName, structureName string
	)

	if rt = reflect.TypeOf(obj); rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	structureName = rt.Name()
	packageName = rt.PkgPath()

	return packageName + string(os.PathSeparator) + structureName
}
