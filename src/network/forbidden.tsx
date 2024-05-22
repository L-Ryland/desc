import { Link } from "react-router-dom"

export default function Forbidden () {
    return (
        <div style={ { color: "#dddddd", fontFamily: "Microsoft Yahei" } }>
            <h1>权限不足！</h1>
            您没有访问该页面的权限！
            <Link to="/关于网站/介绍">返回首页</Link>
        </div>
    )
}
