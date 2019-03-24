package main

// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

var d = make([]int, 10)
var buffer = make([]int, 0)

func append_data(data []byte) {
	file, err := os.OpenFile("data_sonar.txt", os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	defer file.Close()

	len, err := file.WriteString(string(data) + "\n")
	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}
	fmt.Printf("\nLength: %d bytes", len)
	fmt.Printf("\nFile Name: %s", file.Name())
}

func read_data() {
	data, err := ioutil.ReadFile("data_sonar.txt")
	if err != nil {
		log.Panicf("failed reading data from file: %s", err)
	}
	fmt.Printf("\nLength: %d bytes", len(data))
	fmt.Printf("\nData: %s", data)
	fmt.Printf("\nError: %v", err)
}

func get_data(w http.ResponseWriter, r *http.Request) {
	length := len(buffer)
	data := make([][2]int, 0)
	for i := 0; i < length; i++ {
		data = append(data, [2]int{i + 1, buffer[i]})
	}
	json_data, err := json.Marshal(data)
	if err == nil {
		w.Write(json_data)
	}
}

func sonar(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		log.Printf("recv: %s", message)

		if err == nil {
			//append_data(message)
			err = json.Unmarshal(message, &d)
			if err != nil {
				log.Println(err)
			}
			buffer = append(buffer, d...)
			if len(buffer) >= 30 {
				buffer = buffer[len(buffer)-30:]
			}
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, nil)
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/sonar", sonar)
	http.HandleFunc("/", home)
	http.HandleFunc("/getsonar", get_data)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>ECharts</title>
    <!-- 引入 echarts.js -->
	<script src="https://cdnjs.cloudflare.com/ajax/libs/echarts/4.2.1/echarts.min.js"></script>
	<script src="https://cdn.bootcss.com/jquery/3.3.1/jquery.min.js"></script>
</head>
<body>
    <!-- 为ECharts准备一个具备大小（宽高）的Dom -->
    <div id="main" style="width: 600px;height:400px;"></div>
    <script type="text/javascript">
        // 基于准备好的dom，初始化echarts实例
        var myChart = echarts.init(document.getElementById('main'));
        var receive = [];

        // 指定图表的配置项和数据
        option = {
            xAxis: {
                type: 'value',
            },
            yAxis: {
				type: 'value'
            },
            series: [{
                data: [10,2,3,40,5,6,7,8,9,10],
                type: 'line'
            }]
        };


        // 使用刚指定的配置项和数据显示图表。
		myChart.setOption(option);
		
		var secs = 300; //倒计时的秒数 
		function doUpdate(num)   
		{   
			if (num % 3 == 0)
			{
				$.ajax({
					type: 'get', //发送请求类型为POST
					url: 'getsonar', //请求页面的URL，此页面即为上面所述提供JSON数据的页面，传递参数ShowChart，后台解析并调用相应的函数
					data: {},
					dataType: 'json', //请求数据类型为JSON
					async: true, //是否为异步加载，true为异步，false为同步
					success: function (result) { //请求成功：result为后台返回数据
						if (result) {
							option.series[0].data = result;//将得到的数据赋值给option的Series属性
							myChart.setOption(option);
							console.log(result)
							console.log(option)
						}
					},
					error: function () { //请求失败
						alert("Error");
					}
				});
			}  
		}  
		
		for(var i = secs; i >= 0; i--)   
		{   
			window.setTimeout("doUpdate(" + i + ")", (secs-i) * 1000);
		}  
    </script>
</body>
</html>
`))
