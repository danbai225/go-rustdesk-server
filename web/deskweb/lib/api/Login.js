import request  from "../request"

export function getLogin(){
    return request({
        url:"/api/login",
        method:"GET"
    })
}