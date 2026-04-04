import { useId, type InputHTMLAttributes } from "react"
import styles from "./index.module.sass"
import classnames from "classnames"

interface IInput extends InputHTMLAttributes<HTMLInputElement> {
	className?: string
	label?: string
	error?: string
}

export function Input({
	className,
	label,
	type = "text",
	error,
	...props
}: IInput) {
	const htmlFor = useId()
	return (
		<div className={styles.input__wrapper}>
			<label className={styles.label} htmlFor={htmlFor}>
				{label}
			</label>
			<input
				id={htmlFor}
				type={type}
				className={classnames(className, styles.input, {
					[styles.input_error]: error,
				})}
				{...props}
			/>
			{error ? <p className={styles.error}>{error}</p> : null}
		</div>
	)
}
