wrk.method = "POST"

wrk.body  = '{
                 "iid": "1m8loki",
                 "imageUrl": "//s5.mogucdn.com/mlcdn/c45406/180916_4di1ek7k3kha3klk02185678b025d_640x960.jpg",
                 "title": "中长款长袖连衣裙秋装新款2018韩版休闲胖mm大码女装裙子女学生宽松松垮垮中长款卫衣裙外套",
                 "desc": "新款",
                 "price": "54.0",
                 "count": 1
             }'

wrk.headers["Content-Type"] = "application/json"
wrk.headers["Authorization"] = "application/json"

function request()

  return wrk.format('POST', nil, nil, body)

end
