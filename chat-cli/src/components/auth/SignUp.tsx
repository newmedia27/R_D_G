import { Input } from "../ui/form/input/input"
import styles from "./index.module.sass"
import { useForm } from "react-hook-form"
import type { ISignUp } from "../form.types"
import { signUpSchema } from "./schema"
import { yupResolver } from "@hookform/resolvers/yup"
import { signinRequest, signupRequest } from "@/api/auth/auth"
import type { AxiosError } from "axios"
import { setTokenSelector, useAuthStore } from "@/stores/auth"

export default function SignUp() {
	const setToken = useAuthStore(setTokenSelector)

	const {
		register,
		handleSubmit,
		reset,
		formState: { errors },
	} = useForm<ISignUp>({
		defaultValues: {
			email: "newmedia27@gmail.com",
			username: "jart",
			password: "newmediA2",
			passwordConfirm: "newmediA2",
		},
		mode: "onSubmit",
		resolver: yupResolver(signUpSchema),
	})

	async function onSubmit(values: ISignUp): Promise<void> {
		const password = values.password
		const email = values.email
		try {
			await signupRequest(values)
			const { data } = await signinRequest({ email, password })
			setToken(data?.access_token, data.user.id)
			reset()
		} catch (err) {
			const error = err as AxiosError
			console.log("err", error?.response?.data)
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
				{...register("username")}
				label="Username*"
				placeholder="Enter you Username"
				error={errors?.username?.message}
			/>
			<Input
				{...register("password")}
				label="Password*"
				placeholder="Enter you password"
				type="password"
				error={errors?.password?.message}
			/>
			<Input
				{...register("passwordConfirm")}
				label="Password*"
				placeholder="Confirm you password"
				type="password"
				error={errors?.passwordConfirm?.message}
			/>
			<button className={styles.button} type="submit">
				submit
			</button>
		</form>
	)
}
