import type { ModalStep } from ".."
import styles from "./index.module.sass"

type Props = (dirrect: ModalStep) => void

export function NewStep({ setStep }: { setStep: Props }) {
	return (
		<>
			<button
				className={styles.typeBtn}
				onClick={() => setStep("dirrect")}
			>
				<div className={styles.typeIcon}>💬</div>
				<div>
					<div className={styles.typeTitle}>Private chat</div>
					<div className={styles.typeDesc}>
						Один на один з іншим користувачем
					</div>
				</div>
			</button>
			<button className={styles.typeBtn} onClick={() => setStep("group")}>
				<div className={styles.typeIcon}>👥</div>
				<div>
					<div className={styles.typeTitle}>Group chat</div>
					<div className={styles.typeDesc}>
						Додайте кількох учасників
					</div>
				</div>
			</button>
		</>
	)
}
