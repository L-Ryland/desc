import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';
import { Button, ThemeProvider, createTheme } from '@mui/material';
import { Outlet, RouteObject, RouterProvider, createBrowserRouter } from 'react-router-dom';
import NoAuth from './network/noauth';
import Forbidden from './network/forbidden';
import { UserList } from "./pages";
import { UserDataProvider } from './auth';
import { TagProvider } from './data/tag';
import { deepPurple } from '@mui/material/colors';

const theme = createTheme({
    palette: {
        // primary: {
        //     main: '#90caf9',
        // },
        // secondary: deepPurple
    },
});

const webContents = [
    {
        index: true,
        element: <App />,
    },
    {
        path: "noauth",
        element:
            <NoAuth />
    },
    {
        path: "forbidden",
        element:
            <Forbidden />
    },
    {
        path: "users",
        element: <UserList />
    },
    {
        path: ":type",
        element: <App />,
    }
] satisfies RouteObject[];

function browserRouter() {
    return createBrowserRouter([
        {
            path: "/",
            element: <Outlet />,
            children: [...webContents],
        },
    ]);
}

const root = ReactDOM.createRoot(
    document.getElementById('root') as HTMLElement
);
root.render(
    <React.StrictMode>
        <ThemeProvider theme={theme}>
            <UserDataProvider>
                <TagProvider>
                    <RouterProvider router={browserRouter()} />
                </TagProvider>
            </UserDataProvider>
        </ThemeProvider>
    </React.StrictMode>
);
