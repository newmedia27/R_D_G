import type { AxiosResponse } from "axios"
import { api } from ".."
import type { User, UserMap } from "@/components/user.types"

type BatchResponse = Record<string, UserMap>
export async function getUsersByIdsRequest(
	ids: string[]
): Promise<AxiosResponse<BatchResponse>> {
	const url = `/users/batch`
	return await api.post<BatchResponse>(url, { ids })
}

type SearchResponse = Record<string, User[]>
export async function searchUsersRequest(
	search: string
): Promise<AxiosResponse<SearchResponse>> {
	const url = `/users/search?search=${search}`
	return await api.get<SearchResponse>(url)
}

export async function getProfileRequest(): Promise<
	AxiosResponse<Record<string, User>>
> {
	const url = `/users/profile`
	return await api.get<Record<string, User>>(url)
}
