import { Link } from "react-router-dom"
import { Logout } from "../auth"

export default function NoAuth () {
    Logout()
    return (
        <div style={ { color: "#dddddd", fontFamily: "Microsoft Yahei" } }>
            <h1>登录信息错误！</h1>
            您的登录信息已过期，或者所访问的页面权限不足！
            <Link to="/">返回首页</Link>
        </div>
    )
}
