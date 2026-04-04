import { Outlet, useLocation } from "react-router-dom"
import styles from "./index.module.sass"

type Title = "Sign In" | "Sign Up"

export default function AuthLayout() {
	const location = useLocation()
	let title: Title = "Sign Up"
	if (location.pathname === "/auth/signin") {
		title = "Sign In"
	}
	return (
		<section className={styles.main}>
			<div className={styles.form__wrapper}>
				<h1 className={styles.title}>{title}</h1>
				<Outlet />
			</div>
		</section>
	)
}
