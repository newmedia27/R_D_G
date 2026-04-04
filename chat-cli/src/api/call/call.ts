import type { AxiosResponse } from "axios"
import { api } from ".."
import type {
	CallTokenResponse,
	StartCallRequest,
} from "@/components/call.types"

export async function startCallRequest(
	data: StartCallRequest
): Promise<AxiosResponse<CallTokenResponse>> {
	const url = `/calls/token`
	return await api.post<CallTokenResponse>(url, data)
}
