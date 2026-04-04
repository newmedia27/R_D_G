import {
	Children,
	cloneElement,
	useCallback,
	useEffect,
	useRef,
	useState,
	type ReactElement,
	type ReactNode,
} from "react"
import styles from "./index.module.sass"
import { NewStep } from "./new"
import { Dirrect } from "./dirrect"
import { Group } from "./group"
import { useUserStore } from "@/stores/user"

interface TriggerProps {
	onClick?: () => void
}

export type ModalStep = "new" | "dirrect" | "group"

export function AddChatModal({
	children,
}: {
	children: ReactElement<TriggerProps>
}) {
	const [isOpen, setIsOpen] = useState<boolean>(false)
	const [step, setStep] = useState<ModalStep>("new")
	const dialogRef = useRef<HTMLDialogElement>(null)
	const resetSerchedUsers = useUserStore((state) => state.resetSerchedUsers)

	const handleClose = useCallback(() => {
		setIsOpen(false)
		setStep("new")
		resetSerchedUsers()
	}, [resetSerchedUsers])

	useEffect(() => {
		const dialog = dialogRef.current
		if (dialog) {
			if (isOpen) {
				dialog.showModal()
			} else {
				dialog.close()
			}
		}
	}, [isOpen])

	useEffect(() => {
		const dialog = dialogRef.current
		if (!dialog) return

		function handleBackdropClick(e: MouseEvent) {
			if (e.target === dialog) {
				handleClose()
			}
		}

		dialog.addEventListener("click", handleBackdropClick)
		return () => dialog.removeEventListener("click", handleBackdropClick)
	}, [handleClose])

	return (
		<>
			{cloneElement(children, { onClick: () => setIsOpen(true) })}
			<dialog
				ref={dialogRef}
				className={styles.dialog}
				onClose={handleClose}
			>
				<div className={styles.modal}>
					<div className={styles.header}>
						{step !== "new" && (
							<button
								className={styles.backBtn}
								onClick={() => setStep("new")}
							>
								←
							</button>
						)}
						<span className={styles.title}>
							{step === "new" && "New chat"}
							{step === "dirrect" && "Dirrect chat"}
							{step === "group" && "Group chat"}
						</span>
						<button
							className={styles.closeBtn}
							onClick={handleClose}
						>
							✕
						</button>
					</div>
					<ModalStepContent step={step}>
						<NewStep setStep={setStep} />
						<Dirrect onClose={handleClose} />
						<Group onClose={handleClose} />
					</ModalStepContent>
				</div>
			</dialog>
		</>
	)
}
interface ModalContent {
	step: ModalStep
	children: ReactNode
}

const STEP_INDEX: Record<ModalStep, number> = {
	new: 0,
	dirrect: 1,
	group: 2,
}

function ModalStepContent({ step, children }: ModalContent) {
	const childrenToArray = Children.toArray(children)

	return (
		<div className={styles.body}>{childrenToArray[STEP_INDEX[step]]}</div>
	)
}
