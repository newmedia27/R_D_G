import { useEffect, type RefObject } from "react"

export const useClickOutside = (
	target: RefObject<HTMLElement | null>,
	cb: () => void
) => {
	useEffect(() => {
		function handleClick(e: MouseEvent) {
			if (target.current && !target.current.contains(e.target as Node)) {
				cb()
			}
		}
		document.addEventListener("mousedown", handleClick)
		return () => document.removeEventListener("mousedown", handleClick)
	}, [target, cb])
}
