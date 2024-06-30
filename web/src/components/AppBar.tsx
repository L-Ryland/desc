import * as React from 'react';
import { AppBar as MuiAppBar, Box, Toolbar, IconButton, Typography, Menu, Container, Avatar, Button, Tooltip, MenuItem, Link } from "@mui/material";
// import { Menu as MenuIcon } from "@mui/icons-material"
import MenuIcon from '@mui/icons-material/Menu';
import { LoginButton } from "./LoginButton";
import { useCallback, useMemo, useState } from 'react';
import { api } from "../network/api";
import { useNavigate } from 'react-router-dom';

const pages = ['Products', 'Pricing', 'Blog'];
const settings = [
  // 'Profile',
  // 'Account',
  // 'Dashboard',
  'Logout'
];

const checkAuth = async () => {
  try {
    const { data } = await api.get("/auth");
    sessionStorage.setItem("User", JSON.stringify(data));
    return true;
  } catch (error) {
    return false;
  }
};
interface UserData {
  user: string;
  role: string;
}
export const AppBar = () => {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [loggedInUser, setLoggedInUser] = useState<UserData>();
  const [anchorElNav, setAnchorElNav] = React.useState<null | HTMLElement>(null);
  const [anchorElUser, setAnchorElUser] = React.useState<null | HTMLElement>(null);

  const navigate = useNavigate();

  const reactCheckAuth = useCallback(async () => {
    const isLoggedIn = await checkAuth();
    setIsLoggedIn(isLoggedIn);
    if (isLoggedIn) {
      const userData = sessionStorage.getItem('User');
      setLoggedInUser(JSON.parse(userData || '{}'));
    }
  }, []);
  React.useEffect(() => {
    reactCheckAuth();
  }, [reactCheckAuth]);
  const isAdmin = useMemo(() => {
    const _loggedInUser = loggedInUser;
    return !!loggedInUser && loggedInUser.role === 'admin';
  }, [loggedInUser]);

  const handleOpenNavMenu = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorElNav(event.currentTarget);
  };
  const handleOpenUserMenu = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorElUser(event.currentTarget);
  };

  const handleCloseNavMenu = (e) => {
    setAnchorElNav(null);
  };
  const handlePageMenuAction = (page: string) => {
    if (page) {
      navigate(page);
    }
  };

  const handleCloseUserMenu = () => {
  };
  const handleSettingsMenuAction = async (setting: string) => {
    if (setting === 'Logout') {
      await api.post('/logout');
      sessionStorage.removeItem('isLoggedIn');
      sessionStorage.removeItem('User');
      setAnchorElUser(null);
      window.location.reload();
    }
  };

  return (
    <MuiAppBar position="static">
      <Container maxWidth="xl">
        <Toolbar disableGutters>
          <Typography
            variant="h6"
            noWrap
            component="a"
            href="#app-bar-with-responsive-menu"
            sx={{
              mr: 2,
              display: { xs: 'none', md: 'flex' },
              fontFamily: 'monospace',
              fontWeight: 700,
              letterSpacing: '.3rem',
              color: 'inherit',
              textDecoration: 'none',
            }}
          >
            LOGO
          </Typography>

          <Box sx={{ flexGrow: 1, display: { xs: 'flex', md: 'none' } }}>
            <IconButton
              size="large"
              aria-label="account of current user"
              aria-controls="menu-appbar"
              aria-haspopup="true"
              onClick={handleOpenNavMenu}
              color="inherit"
            >
              <MenuIcon />
            </IconButton>
            <Menu
              id="menu-appbar"
              anchorEl={anchorElNav}
              anchorOrigin={{
                vertical: 'bottom',
                horizontal: 'left',
              }}
              keepMounted
              transformOrigin={{
                vertical: 'top',
                horizontal: 'left',
              }}
              open={Boolean(anchorElNav)}
              onClose={handleCloseNavMenu}
              sx={{
                display: { xs: 'block', md: 'none' },
              }}
            >
              {pages.map((page) => (
                <MenuItem key={page} onClick={() => handlePageMenuAction(page)}>
                  <Typography textAlign="center">{page}</Typography>
                </MenuItem>
              ))}
            </Menu>
          </Box>
          <Typography
            variant="h5"
            noWrap
            component="a"
            href="#app-bar-with-responsive-menu"
            sx={{
              mr: 2,
              display: { xs: 'flex', md: 'none' },
              flexGrow: 1,
              fontFamily: 'monospace',
              fontWeight: 700,
              letterSpacing: '.3rem',
              color: 'inherit',
              textDecoration: 'none',
            }}
          >
            LOGO
          </Typography>
          <Box sx={{ flexGrow: 1, display: { xs: 'none', md: 'flex' } }}>
            {pages.map((page) => (
              <Button
                key={page}
                onClick={handleCloseNavMenu}
                sx={{ my: 2, color: 'white', display: 'block' }}
              >
                <Link href={"/" + page} color="inherit" underline="none">{page}</Link>
                {/* {page} */}
              </Button>
            ))}
          </Box>

          <Box sx={{ flexGrow: 0 }}>
            {isLoggedIn ? (
              <>
                {isAdmin && <UserManageButton />}
                <Tooltip title="Open settings">
                  <IconButton onClick={handleOpenUserMenu} sx={{ p: 0 }}>
                    <Avatar alt={loggedInUser?.user} src="/static/images/avatar/2.jpg" />
                  </IconButton>
                </Tooltip>
                <Menu
                  sx={{ mt: '45px' }}
                  id="menu-appbar"
                  anchorEl={anchorElUser}
                  anchorOrigin={{
                    vertical: 'top',
                    horizontal: 'right',
                  }}
                  keepMounted
                  transformOrigin={{
                    vertical: 'top',
                    horizontal: 'right',
                  }}
                  open={Boolean(anchorElUser)}
                  onClose={handleCloseUserMenu}
                >
                  {settings.map((setting) => (
                    <MenuItem key={setting} onClick={() => handleSettingsMenuAction(setting)}>
                      <Typography textAlign="center">{setting}</Typography>
                    </MenuItem>
                  ))}
                </Menu>
              </>
            ) : <LoginButton />}
          </Box>
        </Toolbar>
      </Container>
    </MuiAppBar>
  );
};

const UserManageButton = () => {
  return (
    <Button variant="text" sx={{ mr: 2, color: "white" }}>
      <Link href={"/users"} color="inherit" underline="none">用戶管理</Link>
    </Button>
  );
};
