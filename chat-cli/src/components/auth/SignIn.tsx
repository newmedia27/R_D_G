import { setTokenSelector, useAuthStore } from "@/stores/auth"
import { Input } from "../ui/form/input/input"
import styles from "./index.module.sass"
import { useForm } from "react-hook-form"
import type { ISignIn } from "../form.types"
import { yupResolver } from "@hookform/resolvers/yup"
import { signinSchema } from "./schema"
import type { AxiosError } from "axios"
import { signinRequest } from "@/api/auth/auth"

export default function SignIn() {
	const setToken = useAuthStore(setTokenSelector)
	const {
		register,
		handleSubmit,
		reset,
		formState: { errors },
	} = useForm<ISignIn>({
		defaultValues: {
			email: "newmedia27@gmail.com",
			password: "newmediA2",
		},
		mode: "onSubmit",
		resolver: yupResolver(signinSchema),
	})

	async function onSubmit(values: ISignIn): Promise<void> {
		try {
			const { data } = await signinRequest(values)
			setToken(data?.access_token, data.user.id)
			reset()
		} catch (err) {
			const error = err as AxiosError
			console.error("error", error?.response?.data)
		}
	}

	return (
		<form className={styles.form} onSubmit={handleSubmit(onSubmit)}>
			<Input
				{...register("email")}
				label="Email*"
				placeholder="Enter you email"
				error={errors?.email?.message}
			/>
			<Input
				{...register("password")}
				label="Password*"
				placeholder="Enter you password"
				type="password"
				error={errors?.password?.message}
			/>
			<button className={styles.button} type="submit">
				submit
			</button>
		</form>
	)
}
