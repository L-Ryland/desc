import { Button, Card, CardContent, Chip, IconButton, Link, Paper, TextField, Tooltip, Typography } from "@mui/material"
import { api } from "../network/api"
import { MdDelete, MdModeEdit } from "react-icons/md"
import { ShowTags, TagContext, TagData } from "./tag"
import { useContext, useState } from "react"

export class WebData {
    public ID?: number
    public Name?: string
    public Url!: string // This is a single Primitive
    public Description?: string
    public Tags!: string[] // This is a Primitive Array

    constructor ( url: string, tags: string[], name?: string, description?: string, id?: number ) {
        this.Name = name
        this.Url = url
        this.Description = description
        this.Tags = tags
        this.ID = id
    }

    public async save (): Promise<string> {
        try {
            if ( this.ID ) {
                const response = await api.patch( `/web/` + this.ID, this )
                console.log( response.data )
                return 'success'
            } else {
                const response = await api.post( `/web`, this )
                console.log( response.data )
                return 'success'
            }
        } catch ( error: any ) {
            if ( error.response.status === 409 ) {
                return `${ this.Url } 网址已存在`
            }
            return 'Error add web ' + this.Url + ':' + error
        }
    }

    public delete () {
        console.log( this )
        if ( !this.ID ) {
            return
        }
        api.delete( `/web/` + this.ID )
            .then( response => {
                console.log( response.data )
            } )
            .catch( error => {
                console.log( 'Error delete web ' + this.Url + ':' + error )
            } )
    }

    public show ( addTag: ( tag: string ) => void, deleteData: () => void, edit: () => void ) {
        return (
            <Card key={ this.ID } style={ { margin: '5px' } }>
                <CardContent>
                    <div style={ { display: 'flex', justifyContent: 'space-between', alignItems: 'center' } }>
                        <Tooltip title={
                            this.Description ? (
                                <Typography variant="body1">{ this.Description }</Typography>
                            ) : null
                        } arrow>
                            <Typography variant="h6">{ this.Name }</Typography>
                        </Tooltip>
                        <div>
                            <IconButton onClick={ () => {
                                edit()
                            } }>
                                <MdModeEdit />
                            </IconButton>
                            <IconButton onClick={ () => {
                                this.delete()
                                deleteData()
                            } }>
                                <MdDelete />
                            </IconButton>
                        </div>
                    </div>
                    <Link href={ this.Url } target="_blank" rel="noopener noreferrer">{ this.Url }</Link>

                    <div style={ { display: 'flex', alignItems: 'center' } }>
                        <Typography variant="body2">
                            标签:
                        </Typography>
                        { this.Tags?.map( ( tag, index ) => (
                            <Chip
                                key={ index }
                                label={ tag }
                                onClick={ () => addTag( tag ) }
                                variant="outlined"
                                style={ { margin: '0.5rem' } }
                            />
                        ) ) }
                    </div>
                </CardContent>
            </Card>
        )
    }
}

export async function Search ( tags: string[] ): Promise<WebData[]> {
    try {
        const response = await api.get( `/web/` + tags.join( "," ) )
        console.log( response.data )

        // 创建 WebData 数组
        const webDataArray: WebData[] = response.data.map( ( data: WebData ) => {
            return new WebData( data.Url, data.Tags, data.Name, data.Description, data.ID )
        } )

        return webDataArray
    } catch ( error ) {
        console.error( 'Error search ' + tags.join( "," ) + ':', error )
        return []
    }
}

export function EditWeb ( data: WebData, setData: ( React.Dispatch<React.SetStateAction<WebData>> ), saveData: () => void ) {
    const ctx = useContext( TagContext )
    if ( !ctx ) {
        throw new Error( 'tagsData must be used within a tagsProvider' )
    }

    const handleNameChange = ( event: React.ChangeEvent<HTMLInputElement> ) => {
        setData( ( prevData ) => ( {
            ...prevData,
            Name: event.target.value,
            save: prevData.save,
            delete: prevData.delete,
            show: prevData.show,
        } ) )
    }

    const handleUrlChange = ( event: React.ChangeEvent<HTMLInputElement> ) => {
        setData( ( prevData ) => ( {
            ...prevData,
            Url: event.target.value,
            save: prevData.save,
            delete: prevData.delete,
            show: prevData.show,
        } ) )
    }

    const handleDescriptionChange = ( event: React.ChangeEvent<HTMLInputElement> ) => {
        setData( ( prevData ) => ( {
            ...prevData,
            Description: event.target.value,
            save: prevData.save,
            delete: prevData.delete,
            show: prevData.show,
        } ) )
    }

    const handleTagsChange = ( tags: string[] ) => {
        setData( ( prevData ) => ( {
            ...prevData,
            Tags: tags,
            save: prevData.save,
            delete: prevData.delete,
            show: prevData.show,
        } ) )
    }

    return (
        <Paper>
            <Typography style={ { margin: '10px' } }>网站数据</Typography>
            <TextField
                required
                label="Name"
                value={ data.Name }
                onChange={ handleNameChange }
                style={ { margin: '5px' } }
            />
            <TextField
                required
                label="Url"
                value={ data.Url }
                onChange={ handleUrlChange }
                style={ { margin: '5px' } }
            />
            <TextField
                required
                label="Description"
                value={ data.Description }
                onChange={ handleDescriptionChange }
                style={ { margin: '5px' } }
            />
            <br />
            <ShowTags
                options={ ctx.tags }
                tags={ data.Tags.map( tag => new TagData( tag ) ) }
                setTags={ ( tags ) => ( handleTagsChange( tags.map( tag => tag.Name ) ) ) }
            />
            <Button
                onClick={ saveData }
                variant="outlined"
                disabled={ data.Tags.length === 0 }
                style={ { margin: '5px' } }
            >
                { data.Tags.length > 0 ? "保存" : "先选择标签再保存" }
            </Button>
        </Paper>
    )
}
