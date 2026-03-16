import type { Chat } from "@/components/chat.types"
import styles from "./index.module.sass"
import type { User } from "@/components/user.types"
import { format } from "date-fns"

interface Props {
	chat: Chat
	member?: User
	active?: boolean
	onClick: () => void
}

export const ChatItem = ({ chat, member, active, onClick }: Props) => {
	const avatar =
		member?.username[0]?.toUpperCase() ?? chat.name[0]?.toUpperCase()
	const name = member?.username ?? chat.name
	return (
		<div
			className={`${styles.item} ${active ? styles.active : ""}`}
			onClick={onClick}
		>
			<div className={styles.avatarWrap}>
				<div className={styles.avatar}>{avatar}</div>
				<div className={styles.onlineDot} />
			</div>
			<div className={styles.info}>
				<div className={styles.top}>
					<span className={styles.name}>{name}</span>
					<span className={styles.time}>
						{format(chat.updated_at, "HH:mm")}
					</span>
				</div>
				<span className={styles.lastMsg}>
					{chat.last_message?.text ?? "Немає повідомлень"}
				</span>
			</div>
		</div>
	)
}
