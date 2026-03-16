import type { AuthResponse, ISignIn, ISignUp } from "@/components/form.types"
import { api } from ".."
import type { AxiosResponse } from "axios"

export async function signupRequest(
	data: ISignUp
): Promise<AxiosResponse<AuthResponse>> {
	const url = "/auth/signup"
	return await api.post<AuthResponse>(url, data)
}

export async function signinRequest(
	data: ISignIn
): Promise<AxiosResponse<AuthResponse>> {
	const url = `/auth/signin`
	return api.post<AuthResponse>(url, data)
}

export async function signOutRequest(): Promise<AxiosResponse<void>> {
	const url = `/auth/logout`
	return await api.post<void>(url)
}
