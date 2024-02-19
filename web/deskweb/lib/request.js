import axios from "axios";
import { config } from "next/dist/build/templates/pages";

const instance = axios.create({
    baseURL:process.env.NEXT_PUBLIC_API_URL,
})

instance.interceptors.request.use(
    config =>{
        return config
    },
    err=>{
        return err
    }
)

instance.interceptors,Response.use(
    response=>{
        return response.data
    },
    err=>{
        return Promise.reject(err)
    }
)

export function request(){
    return instance(config)
}