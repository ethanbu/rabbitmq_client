package app

import (
	"MQClient/config"
	"MQClient/logHelper"
	"MQClient/mqListeners"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func Start() {
	done := make(chan struct{}, 1)
	handleSignal(done)
	runtime.GOMAXPROCS(runtime.NumCPU())
	conf := config.GetSystemConfig()
	security := &mqListeners.QueueSecurity{
		Port: conf.RabbitConfig.Port,
		Host: conf.RabbitConfig.Host,
		Password: conf.RabbitConfig.Pwd,
		UserName: conf.RabbitConfig.User,
	}
	listeners := make(map[string]mqListeners.ListenInterface)
	for _, q := range conf.RabbitConfig.Queues {
		var listener mqListeners.ListenInterface
		switch q.Type {
		case "DataImport": //数据导入
			listener = &mqListeners.ImportExcelListener{
					BaseListener: &mqListeners.BaseListener{
					Queue: &mqListeners.Queue{Name: q.Name, Qos: q.Qos},
					QueueSecurity: security,
				},
			}
		case "CommonRequest": //公共请求
			listener = &mqListeners.CommonRequestQueue{
				BaseListener: &mqListeners.BaseListener{
					Queue: &mqListeners.Queue{Name: q.Name, Qos: q.Qos},
					QueueSecurity: security,
				},
			}
		case "AttendanceEvaluation": //考勤评估
			fallthrough
		case "EmployeeOnboarding": //员工入职
			fallthrough
		case "General": //通用请求，与CommonRequest类似
			listener = &mqListeners.GeneralRequestQueue{
				BaseListener: &mqListeners.BaseListener{
					Queue: &mqListeners.Queue{Name: q.Name, Qos: q.Qos},
					QueueSecurity: security,
				},
			}
		case "PayrollCalculation":
			listener = &mqListeners.PayrollCalculationQueue{
				BaseListener : &mqListeners.BaseListener{
					Queue: &mqListeners.Queue{Name: q.Name, Qos: q.Qos},
					QueueSecurity: security,
				},
			}
		}
		if listener != nil {
			go listener.StartListen(done)
			listeners[q.Name] = listener
		}
	}
	fmt.Println("we are listening...")
	<-done
}

func handleSignal(done chan<- struct{}){
	signalChan := make(chan os.Signal, 1)
	//https://www.jianshu.com/p/ae72ad58ecb6
	signal.Notify(signalChan, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGINT)
	go func(){
		sig := <-signalChan
		if sig != nil { //
			logHelper.LogDefault("SYSTEM", "quit")
			close(done)
		}
	}()
}