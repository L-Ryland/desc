import { useEffect, useState } from "react";
import { useSearchParams } from "react-router-dom";

export const useSearchParam = (key: string, defaultVal?: string) => {
  const [search] = useSearchParams();
  const [value, setValue] = useState<string | null>(defaultVal ?? null)

  useEffect(() => {
    let value = search.get(key)
    setValue(value)
  }, [])
  
  return value
}