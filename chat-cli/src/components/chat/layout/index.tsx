import { Outlet } from "react-router-dom"
import styles from "./indes.module.sass"
import { useEffect } from "react"
import { useAuthStore } from "@/stores/auth"
import { useWsStore } from "@/stores/ws"

export default function ChatLayout() {
	useEffect(() => {
		const token = useAuthStore.getState().token
		if (token) {
			useWsStore.getState().connect(token)
		}
		return () => {
			useWsStore.getState().disconnect()
		}
	}, [])
	return (
		<section className={styles.main}>
			<Outlet />
		</section>
	)
}
