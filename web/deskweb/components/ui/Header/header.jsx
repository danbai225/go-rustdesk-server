import style from "./Header.module.css"
import { Input } from "@/components/ui/input"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"

export default function Header(){
    return  <div className={style.HeaderContainer}>
         <div className="flex  justify-between justify-items-center items-center">
         <svg width="28" height="28" viewBox="0 0 28 28" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M0 3.81818C0 1.70946 1.70946 0 3.81818 0H8.90909C11.0178 0 12.7273 1.70946 12.7273 3.81818V12.7273H3.81818C1.70946 12.7273 0 11.0178 0 8.90909V3.81818Z" fill="white"/>
<path d="M28 3.81818C28 1.70946 26.2905 0 24.1818 0H19.0909C16.9822 0 15.2727 1.70946 15.2727 3.81818V12.7273H24.1818C26.2905 12.7273 28 11.0178 28 8.90909V3.81818Z" fill="white"/>
<path d="M0 24.1818C0 26.2905 1.70946 28 3.81818 28H8.90909C11.0178 28 12.7273 26.2905 12.7273 24.1818V15.2727H3.81818C1.70946 15.2727 0 16.9822 0 19.0909V24.1818Z" fill="white"/>
</svg>
         <h2 className={style.title}>DestWeb</h2>
         </div>
         <div className="title">
              <h2 className={style.title}>Welcome Back,Arkhan</h2>
         </div>
          <div >
              <div className={style.search}>
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
<path d="M11.5 21.75C5.85 21.75 1.25 17.15 1.25 11.5C1.25 5.85 5.85 1.25 11.5 1.25C17.15 1.25 21.75 5.85 21.75 11.5C21.75 17.15 17.15 21.75 11.5 21.75ZM11.5 2.75C6.67 2.75 2.75 6.68 2.75 11.5C2.75 16.32 6.67 20.25 11.5 20.25C16.33 20.25 20.25 16.32 20.25 11.5C20.25 6.68 16.33 2.75 11.5 2.75Z" fill="white"/>
<path d="M22.0004 22.7499C21.8104 22.7499 21.6204 22.6799 21.4704 22.5299L19.4704 20.5299C19.1804 20.2399 19.1804 19.7599 19.4704 19.4699C19.7604 19.1799 20.2404 19.1799 20.5304 19.4699L22.5304 21.4699C22.8204 21.7599 22.8204 22.2399 22.5304 22.5299C22.3804 22.6799 22.1904 22.7499 22.0004 22.7499Z" fill="white"/>
</svg>

<Input type="txt" placeholder="Search" className="p-0 border border-gray-300 focus:ring-0 px-4 py-2 focus:outline-nonefocus:ring-0  bg-[#231d29]    rounded-none border-transparent  focus:outline-none  !" />
              </div>

              {/* <i>icon</i> */}
          </div>
          <div className="flex">
                    <h3>Evano</h3>
                    <Avatar>
  <AvatarImage src="https://github.com/shadcn.png" />
  <AvatarFallback>CN</AvatarFallback>
</Avatar>

          </div>
        </div>
}