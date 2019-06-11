package timerTask

//import (
//	"github.com/astaxie/beego/logs"
//	"github.com/robfig/cron"
//)
//
//type CronNew struct {
//	Cr *cron.Cron
//	//Running chan bool
//}
//
//var Cronnew CronNew
//
//func CronInit() {
//	Cronnew.Cr = cron.New()
//	//Cronnew.Running <- false
//}
//
///*
//定时任务1，执行offer下发,阻塞式定时
//*/
//func (c *CronNew) Task1() {
//
//	out := make(chan bool)
//	//spec := "0 */1 * * * ?" //每分钟执行一次
//	//spec := "*/3 * * * * ?" //每三秒执行一次
//	//spec := "*/20 * * * * ?"
//	spec := "* * */12 * * ?" //没12个小时一次
//	go func() {
//		out <- false
//	}()
//
//	_ = c.Cr.AddFunc(spec, func() {
//		select {
//		case b := <-out:
//			if !b {
//				logs.Info("下发offer开始！")
//
//				out <- false
//				//c.Running<-true
//				logs.Info("下发offer结束！")
//			}
//		default:
//			out <- true
//		}
//
//	})
//	c.Cr.Start()
//	//select {}
//}
