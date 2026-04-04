import styles from "./index.module.sass"
import { useState, useRef, useEffect } from "react"

interface Props {
	placeholder?: string
	onSend: (text: string) => void
}

export const MessageInput = ({ placeholder, onSend }: Props) => {
	const [value, setValue] = useState("")
	const textareaRef = useRef<HTMLTextAreaElement>(null)

	useEffect(() => {
		const el = textareaRef.current
		if (!el) return
		el.style.height = "auto"
		el.style.height = `${el.scrollHeight}px`
	}, [value])

	function handleSend() {
		if (!value.trim()) return
		onSend(value)
		setValue("")
	}

	function handleKeyDown(e: React.KeyboardEvent<HTMLTextAreaElement>) {
		if (e.key === "Enter" && !e.shiftKey) {
			e.preventDefault()
			handleSend()
		}
	}

	return (
		<div className={styles.wrapper}>
			<div className={styles.inner}>
				<button className={styles.attachBtn} type="button">
					📎
				</button>
				<textarea
					ref={textareaRef}
					className={styles.input}
					placeholder={placeholder ?? "Message..."}
					value={value}
					rows={1}
					onChange={(e) => setValue(e.target.value)}
					onKeyDown={handleKeyDown}
				/>

				<button
					className={styles.sendBtn}
					type="button"
					disabled={!value.trim()}
					onClick={handleSend}
				>
					↑
				</button>
			</div>
		</div>
	)
}
