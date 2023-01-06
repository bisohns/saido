import { useRef } from "react";

/**
 * @template T
 * @param {T} data
 */
function useDataRef(data: any) {
  const ref = useRef(data);
  ref.current = data;
  return ref;
}

export default useDataRef;
