import { Box, Card, CardHeader, CardContent, CardActionArea, Container, Stack, TextField, Chip } from "@mui/material"
import { api } from "../network/api" 
import { FC, KeyboardEventHandler, SyntheticEvent, useCallback, useEffect, useMemo, useState } from "react";
import { css } from "@emotion/css";
import { TagData } from "../data/tag";
interface CategoryData {
  id: string
  name: string
  tags: TagData[]
}
interface ICategoryListProps {
  onDrop(id: string): void
}
export const CategoryList: FC<ICategoryListProps> = ({ onDrop }) => {
  const [categories, setCategories] = useState<CategoryData[]>([]);

  const getCategories = useCallback(async () => {
    const { data } = await api.get<CategoryData[]>("/categories")
    console.log("data", data)
    setCategories(data)
  }, [])
  useEffect(() => {
    getCategories()
  }, [])
  return ( 
    <Stack direction="row" spacing={2}>
      {categories.map(item => <CategoryCard key={item.id} data={item} successCb={getCategories} onDrop={onDrop} />)}
    </Stack>
  );
}

interface ICategoryCardProps {
  data: CategoryData
  successCb?(): void
  onDrop(id: string): void
}
const CategoryCard: FC<ICategoryCardProps> = ({ data, successCb, onDrop }) => {
  const [isEditingName, setIsEditingName] = useState(false);
  const [newCategoryName, setNewCategoryName] = useState(data.name);

  const tags = useMemo(() => data.tags?.map(item => new TagData(item.Name)) ?? [], [data])

  const updateCategory = useCallback((id: string, data: Partial<CategoryData>) => {
    return api.patch(`/categories/${id}`, data)
  }, [])
  const onNameInputKeyDown = useCallback<KeyboardEventHandler>(async (e) => {
    if (e.key === "Enter") {
      await updateCategory(data.id, { name: newCategoryName })
      setIsEditingName(false)
      successCb?.()
    }
  }, [newCategoryName])

  const removeFromCategory = useCallback(async (e: SyntheticEvent, tag: TagData) => {
    e.stopPropagation()
    await tag.removeFromCategory();
    debugger
    window.location.reload()
  }, [])
  return (
    <Card sx={{ minHeight: 300, width: "100%" }} variant="outlined">
      <CardHeader title={isEditingName 
        ? <TextField label="Name" value={newCategoryName} onChange={e => setNewCategoryName(e.target.value)} autoFocus onKeyDown={onNameInputKeyDown} /> 
        : <div onClick={() => setIsEditingName(true)}
        >{
        data.name}
        </div>} />
        <CardContent sx={{ height: "100%" }} onDragOver={e => e.preventDefault()} 
          onDrop={() => onDrop(data.id)}>
            {tags.map((item, index) => <Chip key={index} label={item.Name} clickable={false} onDelete={(e) => removeFromCategory(e, item)} />)}
          </CardContent>
    </Card>
  )
}