import React from 'react'
import ReactDOM from 'react-dom/client'
import './index.css'
import App from './App'
import { Button, ThemeProvider, createTheme } from '@mui/material'
import { RouterProvider, createBrowserRouter } from 'react-router-dom'
import NoAuth from './network/noauth'
import Forbidden from './network/forbidden'
import { UserDataProvider } from './auth'
import { TagProvider } from './data/tag'
import { deepPurple } from '@mui/material/colors'

const theme = createTheme( {
    palette: {
        // primary: {
        //     main: '#90caf9',
        // },
        // secondary: deepPurple
    },
} )

const webContents = [
    {
        path: "noauth",
        element:
            <NoAuth />
    },
    {
        path: "forbidden",
        element:
            <Forbidden />
    }
]

function browserRouter () {
    return createBrowserRouter( [
        {
            path: "*",
            element: <App />,
            children: [ ...webContents ],
        }
    ] )
}

const root = ReactDOM.createRoot(
    document.getElementById( 'root' ) as HTMLElement
)
root.render(
    <React.StrictMode>
        <ThemeProvider theme={ theme }>
            <UserDataProvider>
                <TagProvider>
                    <RouterProvider router={ browserRouter() } />
                </TagProvider>
            </UserDataProvider>
        </ThemeProvider>
    </React.StrictMode>
)
