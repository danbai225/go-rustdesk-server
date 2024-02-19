"use client"
import React, { useState,useEffect } from 'react';
import style from "./Login.module.css"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { 
  Card ,
  CardHeader,
  CardTitle,
  CardContent,
  CardFooter
} from "@/components/ui/card"


export default function Login() {

  const [clicked, setclicked] = useState(false);
  const [ColorBoolean,setBoolean] = useState(false);
  // 账号
  const [inputValue, setInputValue] = useState('');
//  密码
const [inputPassword,setPassword] = useState("");
  // 定义处理输入框变化的事件处理程序
  const handleInputChange = (event) => {
    // 更新状态变量的值
    setInputValue(event.target.value);
  };
  const handlePasswordChange = (event) =>{
     setPassword(event.target.value);
  }
  const handleSubmit = (event) => {
    setclicked(!clicked);
    // 阻止默认表单提交行为
    event.preventDefault();
    // 在这里可以使用 inputValues 进行其他操作
    if(inputValue || inputPassword){
       setBoolean(false)
      //  登陆之后操作
      
    }else{
      setBoolean(true)
    }

  };
  return (
    <div className={style.Container_div}>
      <div className={style.colustm_right}>
      <svg width="451" height="352" viewBox="0 0 451 352" fill="none" xmlns="http://www.w3.org/2000/svg">
<path fillRule="evenodd" clipRule="evenodd" d="M2.74592 -240.125C6.45309 -275.276 13.6302 -310.131 24.2774 -344.134L374.866 -694.723C408.869 -705.37 443.724 -712.547 478.875 -716.254L2.74592 -240.125ZM167.745 205.888C182.089 219.441 196.998 232.05 212.394 243.717L962.717 -506.606C951.05 -522.002 938.441 -536.911 924.888 -551.255L167.745 205.888ZM321.784 307.7C303.114 299.59 284.812 290.366 266.979 280.027L999.027 -452.021C1009.37 -434.188 1018.59 -415.886 1026.7 -397.216L321.784 307.7ZM388.76 331.618C411.114 337.958 433.8 342.816 456.662 346.194L1065.19 -262.338C1061.82 -285.2 1056.96 -307.886 1050.62 -330.24L388.76 331.618ZM633.215 343.015C603.014 348.597 572.427 351.572 541.814 351.937L1070.94 -177.186C1070.57 -146.573 1067.6 -115.986 1062.01 -85.7854L633.215 343.015ZM768.272 298.852C820.978 273.408 870.379 238.825 914.102 195.102C957.825 151.379 992.408 101.978 1017.85 49.2719L768.272 298.852ZM878.045 -595.306L123.694 159.045C110.977 143.803 99.239 128.021 88.4798 111.781L830.781 -630.52C847.021 -619.761 862.803 -608.023 878.045 -595.306ZM55.2624 54.104L773.104 -663.738C754.193 -673.104 734.844 -681.285 715.167 -688.279L30.7207 -3.83258C37.7155 15.8438 45.8961 35.193 55.2624 54.104ZM10.8703 -74.8765C5.92033 -98.863 2.62668 -123.123 0.989391 -147.474L571.526 -718.011C595.877 -716.373 620.137 -713.08 644.124 -708.13L10.8703 -74.8765Z" fill="#4044ED"/>
</svg>
      </div>
 
   <div  className={style.colustm_left}>
   <svg width="329" height="814" viewBox="0 0 329 914" fill="none" xmlns="http://www.w3.org/2000/svg">
<circle cx="-128" cy="457" r="457" fill="#4044ED"/>
<circle cx="38" cy="38" r="38" transform="matrix(1 0 0 -1 118 141)" fill="#93DFFF"/>
</svg>

   </div>
  <div className={style.colustm_bottom}>
  <svg width="480" height="432" viewBox="0 0 480 432" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M120.365 117.539L400.913 196.295L192.434 399.879L120.365 117.539Z" fill="#4044ED" fillOpacity="0.5"/>
<path d="M225.502 211.606L506.05 290.363L297.571 493.946L225.502 211.606Z" fill="#4044ED" fillOpacity="0.5"/>
<path d="M334.603 304.432L615.15 383.188L406.671 586.771L334.603 304.432Z" fill="#4044ED" fillOpacity="0.5"/>
</svg>

  </div>

   <Card className="w-[350px] font  border-0 ">
   <form onSubmit={handleSubmit} >
      <CardHeader>
        <CardTitle  className="font-login text-left     font-normal text-4xl">LOGIN</CardTitle>
      </CardHeader>

      <CardContent>
      
          <div className="grid w-full items-center gap-4 text-black  ">
            <div className="flex flex-col space-y-1.5">
              <Input id="name" placeholder="Username"  value={inputValue} onChange={handleInputChange}    className={ColorBoolean?'border-2 border-rose-600':''} />
            </div>
            <div className="flex flex1-col space-y-1.5">
            <Input id="name"  value={inputPassword} onChange={ handlePasswordChange}  placeholder="Password"  className={ColorBoolean?'border-2 border-rose-600':''}  />
            </div>
          </div>
    
      </CardContent>
      <CardFooter className="flex justify-between">

        <Button  type="submit" variant={clicked?'outline':''}  className="w-full bg-[#4044ED]">Login</Button>

      </CardFooter>
      </form>
    </Card>

    </div>
 
  )
}
