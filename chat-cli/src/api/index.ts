import axios, {
	AxiosError,
	type AxiosResponse,
	type InternalAxiosRequestConfig,
} from "axios"
import { useAuthStore } from "../stores/auth"
import { routeNavigate as navigate } from "@/routes"

const baseUrl = import.meta.env.VITE_API_URL

interface CustomInternalAxiosRequestConfig extends InternalAxiosRequestConfig {
	startTime?: number
}

interface RefreshResponse {
	access_token: string
}

interface AuthStore {
	state: {
		token: string | null
	}
}

const api = axios.create({
	withCredentials: true,
	baseURL: baseUrl,
})

let refreshPromise: Promise<AxiosResponse<RefreshResponse>> | null = null

export const refreshToken = async (token: string): Promise<string | null> => {
	if (!refreshPromise) {
		refreshPromise = axios.post<RefreshResponse>(
			`${baseUrl}/auth/refresh`,
			null,
			{
				headers: { authorization: `Bearer ${token}` },
				withCredentials: true,
			}
		)
	}

	try {
		const { data } = await refreshPromise
		useAuthStore.getState().setTokenAfterRefresh(data.access_token)
		return data.access_token
	} catch (error) {
		console.log((error as AxiosError).response)
		useAuthStore.getState().resetToken()
		return null
	} finally {
		refreshPromise = null
	}
}

api.interceptors.request.use(
	async (config: CustomInternalAxiosRequestConfig) => {
		if (refreshPromise) {
			await refreshPromise
		}
		const storage = localStorage.getItem("auth-storage") || "{}"

		const auth: AuthStore = JSON.parse(storage)
		let token = auth?.state?.token || null
		if (token) {
			const payload = JSON.parse(atob(token.split(".")[1]))
			const now = Date.now() / 1000
			if (payload.exp - now < 60) {
				const newToken = await refreshToken(token)
				if (newToken) {
					token = newToken
				}
			}
			config.headers.Authorization = `Bearer ${token}`
			config["startTime"] = new Date().valueOf()
		} else {
			delete config.headers.Authorization
		}
		return config
	}
)

api.interceptors.response.use(
	(response: AxiosResponse) => response,
	async (error: AxiosError) => {
		if (!error.response) {
			return Promise.reject(error)
		}
		if (
			error.request.responseType === "blob" &&
			error.response?.data instanceof Blob &&
			error.response.data.type?.toLowerCase().includes("json")
		) {
			try {
				const text = await error.response.data.text()
				error.response.data = JSON.parse(text)
			} catch (e) {
				console.log("e", e)
			}
		}

		const status = error.response?.status as number | undefined
		switch (status) {
			case 401:
				useAuthStore.getState().resetToken()
				navigate("/auth/signin", { replace: true })
				break

			case 403:
				console.error("Access denied")
				break

			case 404:
				console.error("Resource not found")
				break

			case 500:
				console.error("Server error")
				break
		}
		throw error
	}
)

export { api }
