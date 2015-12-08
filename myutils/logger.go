package myutils

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	htmlLogger   *log.Logger
	MyLogger     *log.Logger
	htmlLogFile  *os.File
	xlongLogFile *os.File
)

var currentDay int

const (
	TimeLayout   = "2006-01-02 15:04:05"
	YearMonthDay = "2006-01-02"
)

const MaxLogFileSize = 1024 * 1024 * 20

// html log color
const (
	CLR_UPDATE  = "Slateblue"
	CLR_PULL    = "Chartreuse"
	CLR_QUERY   = "Coral"
	CLR_INSPECT = "Turquoise"
	CLR_ERROR   = "Fuchsia"
)

// html log class
const (
	CLS_UPDATE  = "update"
	CLS_PULL    = "pull"
	CLS_QUERY   = "query"
	CLS_INSPECT = "inspect"
	CLS_ERROR   = "error"
)

func LogBanner() {
	htmlLogger.Println("<img src=res/pix4banner.png alt='Logging Delicious' title='Log is delicious, ...MmKay?  -- Mr. Mackey, Jr.'/>")
}

func LogButtons() {
	htmlLogger.Println(`
		<table>
		<tr>
		<td>
		<img id=buttonAll src=res/pix4universe.png onclick="show('all')" />
		<img id=buttonUpdate src=res/pix4slateblue.png onclick="show('update')" />
		<img id=buttonPull src=res/pix4chartreuse.png onclick="show('pull')" />
		<img id=buttonQuery src=res/pix4coral.png onclick="show('query')" />
		<img id=buttonInspect src=res/pix4turquoise.png onclick="show('inspect')" />
		<img id=buttonError src=res/pix4fuchsia.png onclick="show('error')" />`)
}

func LogWidget() {
	htmlLogger.Println(`
		<div>
		<select id=typeSel onchange="show(this.value)">
		<option value="all">All</option>
		<option value="update">Update</option>
		<option value="pull">Pull</option>
		<option value="query">Query</option>
		<option value="inspect">Inspect</option>
		<option value="error">Error</option>
		</select>
		</div>`)
}

func LogScript() {
	htmlLogger.Println(`
		<script type='text/javascript'>
		function show(val) {
			document.getElementById('typeSel').value = val;
			var allItems = document.getElementsByClassName('filter');
			if(val == 'all') {
				for(var i = 0, length = allItems.length; i < length; i++) {
					if(allItems[i].style.display == 'none') {
						allItems[i].style.display = '';
					}
				}
			} else {
				for(var i = 0, length = allItems.length; i < length; i++) {
					if(allItems[i].style.display == '') {
						allItems[i].style.display = 'none';
					}
				}
				var selectedItems = document.getElementsByClassName(val);
				for(var i = 0, length = selectedItems.length; i < length; i++) {
					if(selectedItems[i].style.display == 'none') {
						selectedItems[i].style.display = '';
					}
				}
			}
		}
		</script>`)
}

func LogPlain(msg string) {
	htmlLogger.Println("<p>" + time.Now().Format(TimeLayout) + ": " + msg + "</p>")
}

func LogBold(msg string) {
	htmlLogger.Println("<p><b>" + time.Now().Format(TimeLayout) + ": " + msg + "</b></p>")
}

func LogAlert(msg string) {
	htmlLogger.Println(`<p class="filter error" style=display:><font color=Fuchsia>` + time.Now().Format(TimeLayout) + ": " + msg + "</font></p>")
}

func LogColor(msg string, clr string) {
	htmlLogger.Println("<p><font color=" + clr + ">" + time.Now().Format(TimeLayout) + ": " + msg + "</font></p>")
}

func LogClass(msg string, cls string) {
	htmlLogger.Println(`<p class="filter ` + cls + `"` + " style=display:>" + time.Now().Format(TimeLayout) + ": " + msg + "</p>")
}

func LogColorClass(msg string, clr string, cls string) {
	htmlLogger.Println(`<p class="filter ` + cls + `"` + " style=display:><font color=" + clr + ">" + time.Now().Format(TimeLayout) + ": " + msg + "</font></p>")
}

func LogXlong(msg string) {
	MyLogger.Println(msg)
}

func SetupLogs() {
	var err error
	ok, _ := Exists("log")
	if !ok {
		err = os.Mkdir("log", os.ModeDir|0755)
		if err != nil {
			fmt.Println("Counld not create log directory")
			fmt.Println("error: ", err)
		}
	}

	//CreatePatchLog()
	CreateXlongLog()
}

func CreatePatchLog() {
	var err error
	htmlLogFile, err = os.OpenFile("log/log"+time.Now().Format(YearMonthDay)+".html", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Could not open/create patch log file")
		fmt.Println("error: ", err)
	}
	htmlLogger = log.New(htmlLogFile, "", 0)
	patchFileInfo, err := htmlLogFile.Stat()
	if err != nil {
		fmt.Println("error: ", err)
	}
	if patchFileInfo.Size() == 0 {
		LogScript()
		LogBanner()
		LogButtons()
		LogWidget()
	}
}

func CreateXlongLog() {
	var err error
	xlongLogFile, err = os.OpenFile("log/log"+time.Now().Format(YearMonthDay), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Counld not open/create xlong log file")
		fmt.Println("error: ", err)
	}
	MyLogger = log.New(xlongLogFile, "", log.LstdFlags|log.Lshortfile)
}

func init() {
	currentDay = time.Now().Day()
	SetupLogs()

	go func() {
		ticker := time.Tick(15 * time.Minute)
		for _ = range ticker {
			//patchFileInfo, err := htmlLogFile.Stat()
			//if err != nil {
			//fmt.Println("error: ", err)
			//}
			//if patchFileInfo.Size() > MaxLogFileSize {
			//CreatePatchLog()
			//}

			if currentDay != time.Now().Day() {
				currentDay = time.Now().Day()
				CreateXlongLog()
			}
		}
	}()
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func DeInit() {
	var err error
	if htmlLogFile != nil {
		MyLogger.Println("Closing HTML logger file.")
		err = htmlLogFile.Close()
		if err != nil {
			MyLogger.Println(err.Error())
		}
	}
	if xlongLogFile != nil {
		MyLogger.Println("Closing logger file, no more logs from here on.")
		err = xlongLogFile.Close()
		if err != nil {
			MyLogger.Println(err.Error())
		}
	}
}
