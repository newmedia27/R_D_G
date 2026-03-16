import classNames from "classnames"
import styles from "./index.module.sass"

interface Props {
	avatar?: string
	email: string
	username: string
	selected?: boolean
	onClick: () => void
}

export const UserItem = ({
	avatar,
	username,
	email,
	selected,
	onClick,
}: Props) => {
	return (
		<div
			className={classNames(styles.item, { [styles.selected]: selected })}
			onClick={onClick}
		>
			<div className={styles.avatar}>{avatar}</div>
			<div className={styles.info}>
				<span className={styles.username}>{username}</span>
				<span className={styles.email}>{email}</span>
			</div>
			{selected && <div className={styles.checkmark}>✓</div>}
		</div>
	)
}
