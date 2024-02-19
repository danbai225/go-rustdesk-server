import Header from "@/components/ui/Header/header"
import style from "./Home.module.css"
import Layout from "@/components/ui/layout"
export default function Home() {
    return  <> <Layout>
      <div className={style.backgroundColor}>
      <Header/>
      </div>
 </Layout>
</>
  }