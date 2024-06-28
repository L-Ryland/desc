import React, { useState, useEffect, createContext, FC, ReactNode, useContext } from 'react'
import { api } from './network/api'

import { Navigate } from "react-router-dom"
import Cookies from "js-cookie"

export interface UserData {
    name: string
    heros: Record<number, boolean>
    role: number
}

export interface PrivateRouteProps {
    auth: string
    children: JSX.Element
}

export function PrivateRoute ( { auth, children }: PrivateRouteProps ) {
    //const token = localStorage.getItem( auth )
    let token = Cookies.get( auth )
    return token ? children : <Navigate to="/关于网站/介绍" />
}


interface UserDataContextProps {
    userData: UserData | null
    setUserName: ( newUserName: string | null ) => void
    setUserData: ( userData: UserData ) => void
}

export const UserDataContext = createContext<UserDataContextProps | undefined>( undefined )

export const UserDataProvider: FC<{ children: ReactNode }> = ( { children } ) => {
    const [ userData, setUserData ] = useState<UserData | null>( () => {
        const storedUserData = localStorage.getItem( "userData" )
        return storedUserData ? JSON.parse( storedUserData ) : null
    } )

    const [ userName, setUserNameInternal ] = useState<string | null>( () => Cookies.get( "user" ) || null )


    const setUserName = ( newUserName: string | null ) => {
        setUserNameInternal( newUserName )
    }

    useEffect( () => {
        if ( userName !== null ) {
            Cookies.set( "user", userName, { sameSite: "Strict" } )
        } else {
            Cookies.remove( "user" )
            setUserData( null )
        }
    }, [ userName ] )


    // 将 UserData 存储到 localStorage
    useEffect( () => {
        if ( userData !== null ) {
            localStorage.setItem( "userData", JSON.stringify( userData ) )
        } else {
            localStorage.removeItem( 'userData' )
        }
    }, [ userData ] )

    return (
        <UserDataContext.Provider value={ { userData, setUserName, setUserData } }>
            { children }
        </UserDataContext.Provider>
    )
}

export function Logout () {
    const ctx = useContext( UserDataContext )
    if ( !ctx ) {
        throw new Error( 'useUserData must be used within a UserDataProvider' )
    }
    ctx.setUserName( null )
    Cookies.remove( "session" )
}
