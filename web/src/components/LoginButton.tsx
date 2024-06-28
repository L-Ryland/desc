import { Box, Button, Modal, TextField, Typography } from "@mui/material"
import { LoadingButton } from "@mui/lab"
import { api } from "../network/api"
import { useCallback, useState } from "react";

const dialogStyle = {
  position: 'absolute' as 'absolute',
  top: '50%',
  left: '50%',
  transform: 'translate(-50%, -50%)',
  width: 400,
  bgcolor: 'background.paper',
  border: '2px solid #000',
  boxShadow: 24,
  p: 4,
  display: "flex",
  flexDirection: "column",  
};
export const LoginButton = () => {
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const [username, setUsername] = useState<string>();
  const [password, setPassword] = useState<string>();

  const handleLogin = useCallback(async () => {
    if (!username || !password) return
    setLoading(true)
    try {
      const { data } = await api.post('/login', { Username: username, Password: password })
      localStorage.setItem("User", JSON.stringify(data))
      localStorage.setItem("isLoggedIn", "true")
      window.location.reload()
    } finally{
      setLoading(false) 
      setOpen(false)
    }
  }, [username, password])
  
  return ( 
    <>
      <Button variant="contained" color="secondary" onClick={() => setOpen(true)}>Login</Button>
      <Modal open={open}>
        <Box sx={dialogStyle}>
          <Typography sx={{ marginBottom: 4 }}>Login</Typography>
          <TextField sx={{ marginBottom: 4 }} value={username} label="Username" fullWidth required onChange={e => setUsername(e.target.value)} />
          <TextField sx={{ marginBottom: 4 }} label="Password" type="password" value={password} fullWidth required onChange={e => setPassword(e.target.value)} />
          <Box display="flex" justifyContent="center">
            <LoadingButton variant="contained" color="primary" sx={{ marginRight: 2 }} loading={loading} onClick={handleLogin}>Login</LoadingButton>
            <Button variant="outlined" onClick={() => setOpen(false)}>Cancel</Button>
          </Box>
        </Box>
      </Modal>
    </>
  );
}