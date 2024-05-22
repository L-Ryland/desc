import { FC, PropsWithChildren } from "react";

interface IDraggableElementProps {
    className?: string
}
export const DraggableElement: FC<PropsWithChildren<IDraggableElementProps>> = ({ children, className}) => {
    return (
        <div className={className} draggable>{children}</div>
    )
}