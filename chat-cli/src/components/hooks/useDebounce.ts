import { useEffect, useState } from "react"

export function useDebounce<T>(value: T, delay: number = 1000): T {
	const [debounced, setDebounced] = useState(value)
	useEffect(() => {
		const timer = setTimeout(() => {
			setDebounced(value)
		}, delay)
		return () => clearTimeout(timer)
	}, [delay, value])

	return debounced
}
