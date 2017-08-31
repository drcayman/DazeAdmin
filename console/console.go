package console

import (
	"fmt"
	"bufio"
	"os"
	"DazeProxy/util"
	"strings"
	"github.com/crabkun/DazeAdmin/database"
	"time"
)

func ShowMenu(){
	fmt.Println("**********命令列表**********")
	fmt.Println("help 显示此帮助")
	fmt.Println("users 显示所有用户")
	fmt.Println("add 添加用户")
	fmt.Println("edit 修改用户")
	fmt.Println("del 删除一个用户（比如del 4就是删除掉ID为4的用户）")
	fmt.Println("exit 退出用户管理")
	fmt.Println("****************************")
}
func Start(){
	ShowMenu()
	r:=bufio.NewReader(os.Stdin)
	command:=""
	for{
		fmt.Print(">>>>>>")
		buf,_,err:=r.ReadLine()
		if err!=nil{
			return
		}
		bufstr:=util.B2s(buf)
		fmt.Sscanf(bufstr,"%s",&command)
		switch strings.ToLower(command) {
		case "help":
			ShowMenu()
		case "users":
			users:=database.GetAllUser()
			fmt.Printf("一共有%d个用户\n",database.GetUserCount())
			for _,v:=range users{
				fmt.Printf("ID:%d\t用户名：%s\t过期时间：%s\t所属服务器组：%s\n",v.Id,v.Username,exptimestring(v.Expired),DBgroupToString(v.Group))
			}
		case "add":
			Add()
		case "edit":
			Edit()
		case "del":
			var id int
			n,_:=fmt.Sscanf(bufstr,"%s%d",&command,&id)
			if n!=2{
				fmt.Println("命令格式错误")
				continue
			}
			if flag,name:=database.DeleteById(id);flag{
				fmt.Printf("删除用户成功（ID：%d,用户名：%s）",id,name)
			}else{
				fmt.Println("删除用户失败，或许是ID错误了？")
			}
		case "exit":
			return
		default:
			fmt.Println("命令格式错误，请输入help来查看帮助")
		}
		command=""
	}
}
func Add(){
	var username string
	var password string
	var exptime time.Time
	var group []string
	fmt.Printf("请输入要新建的用户名(不允许空格，留空为退出)：")
	if n,_:=fmt.Scanf("%s",&username);n==0{
		fmt.Println("退出添加向导")
		return
	}
	fmt.Printf("请输入此用户的新密码(不允许空格，留空为退出)：")
	if n,_:=fmt.Scanf("%s",&password);n==0{
		fmt.Println("退出添加向导")
		return
	}
	times:=0
	for{
		fmt.Printf("请输入此用户的过期时间（格式2006-01-02 15:04:05，输入0为永不过期，3次错误为退出）：")
		if buf,_,err:=bufio.NewReader(os.Stdin).ReadLine();err==nil{
			if string(buf)=="0"{
				exptime=time.Time{}
				break
			}
			exptime,err=time.Parse("2006-01-02 15:04:05",string(buf))
			if err!=nil{
				fmt.Println("格式不正确！请重新输入。")
				times++
				if times==3{
					fmt.Println("多次错误，退出添加向导")
					return
				}
				continue
			}
			break
		}
	}
	fmt.Printf("请输入此用户所属服务器组(多个组用空格分开，留空为所有组！)：")
	if buf,_,err:=bufio.NewReader(os.Stdin).ReadLine();err==nil && string(buf)!=""{
		group=strings.Split(string(buf)," ")
	}
	fmt.Printf("新用户信息如下：\n用户名:%s\n密码:%s\n过期时间:%s\n所属服务器组:%s\n",username,password,exptimestring(exptime),groupstringToHuman(group))
	err:=database.AddUser(database.User{
		Username:username,
		Password:util.GetDoubleMd5(password),
		Expired:exptime,
		Group:groupstringToDB(group),
	})
	if err!=nil{
		fmt.Println("新建此用户失败，很可能是用户名已存在！错误代码：",err.Error())
	}else{
		fmt.Println("新建此用户成功！")
	}
}
func exptimestring(t time.Time)string{
	if t.IsZero(){
		return "永不过期"
	}
	return t.Format("2006-01-02 15:04:05")
}
func groupstringToHuman(s []string) string{
	if len(s)==0{
		return "所有组"
	}else{
		return strings.Join(s,",")
	}
}
func groupstringToDB(s []string) string{
	if len(s)==0{
		return ""
	}else{
		for i,_:=range s{
			s[i]="|"+s[i]+"|"
		}
		return strings.Join(s,",")
	}
}
func DBgroupToString(s string) string{
	if len(s)==0{
		return "所有组"
	}else{
		return strings.Replace(s,"|","",-1)
	}
}
func Edit(){
	var id int
	var u database.User
	var flag bool
	var sel int
	fmt.Println("请输出要编辑的用户ID（留空为退出）：")
	if n,_:=fmt.Scanf("%d",&id);n==0{
		fmt.Println("退出编辑向导")
		return
	}
	if u,flag=database.GetUserById(id);!flag{
		fmt.Println("此用户不存在，退出编辑向导！")
		return
	}
	fmt.Printf("此用户信息如下：\n用户名:%s\n过期时间:%s\n所属服务器组:%s\n",
	u.Username,exptimestring(u.Expired),DBgroupToString(u.Group))
	fmt.Println("\n\n要编辑的项目：\n1.用户名\n2.密码\n3.过期时间\n4.所属服务器组\n请输入要编辑的序号（留空为退出）：")
	if n,_:=fmt.Scanf("%d",&sel);n==0{
		fmt.Println("退出编辑向导")
		return
	}
	switch sel {
	case 1:editUserName(id,u)
	case 2:editPassword(id,u)
	case 3:editExptime(id,u)
	case 4:editGroup(id,u)
	default:
		fmt.Println("错误输入，退出编辑向导")
		return
	}
}
func editUserName(id int,u database.User){
	var username string
	fmt.Printf("请输入新用户名(不允许空格，留空为退出)：")
	if n,_:=fmt.Scanf("%s",&username);n==0{
		fmt.Println("退出编辑向导")
		return
	}
	u.Username=username
	database.EditUserById(id,u)
}
func editPassword(id int,u database.User){
	var password string
	fmt.Printf("请输入新密码(不允许空格，留空为退出)：")
	if n,_:=fmt.Scanf("%s",&password);n==0{
		fmt.Println("退出编辑向导")
		return
	}
	u.Password=util.GetDoubleMd5(password)
	database.EditUserById(id,u)
}
func editExptime(id int,u database.User){
	var exptime time.Time
	times:=0
	for{
		fmt.Printf("请输入此用户的过期时间（格式2006-01-02 15:04:05，输入0为永不过期，3次错误为退出）：")
		if buf,_,err:=bufio.NewReader(os.Stdin).ReadLine();err==nil{
			if string(buf)=="0"{
				exptime=time.Time{}
				break
			}
			exptime,err=time.Parse("2006-01-02 15:04:05",string(buf))
			if err!=nil{
				fmt.Println("格式不正确！请重新输入。")
				times++
				if times==3{
					fmt.Println("多次错误，退出编辑向导")
					return
				}
				continue
			}
			break
		}
	}
	u.Expired=exptime
	database.EditUserById(id,u)
}
func editGroup(id int,u database.User){
	var group []string
	fmt.Printf("请输入此用户所属服务器组(多个组用空格分开，留空为所有组！)：")
	if buf,_,err:=bufio.NewReader(os.Stdin).ReadLine();err==nil && string(buf)!=""{
		group=strings.Split(string(buf)," ")
	}
	u.Group=groupstringToDB(group)
	database.EditUserById(id,u)
}