import { Box, Button, Container, Modal, Stack } from "@mui/material";
import { DataGrid, useGridApiRef, type GridColDef } from "@mui/x-data-grid";
import { FC, MutableRefObject, useCallback, useEffect, useState } from "react";
import { api } from "../network/api";
import { FormRadioGroup, FormTextField } from "../components";
import { Form, useForm } from "react-hook-form";
import { GridApiCommunity } from "@mui/x-data-grid/internals";

// const (
// 	RolePlayer RoleLevel = iota
// 	RoleManager
// 	RoleAdmin
// )
enum UserRole {
  RolePlayer,
  RoleManager,
  RoleAdmin
}
interface UserDataRaw {
  id: string;
  Name: string;
  Role: string;
  Heros: string[];
}
const getUsers = async () => {
  const { data } = await api.get<UserDataRaw[]>("/user");
  return data;
};
const columns = [
  { field: "id", headerName: "ID", flex: 1 },
  { field: "Name", align: "center", headerAlign: "center", flex: 1 },
  { field: "Role", align: "center", headerAlign: "center", flex: 1 }
] satisfies GridColDef[];

const roleOptions = [
  { value: UserRole.RolePlayer, label: "RolePlayer" },
  { value: UserRole.RoleManager, label: "RoleManager" },
  { value: UserRole.RoleAdmin, label: "RoleAdmin" },
];

const UserList = () => {
  const [rows, setRows] = useState<UserDataRaw[]>([]);
  const getUsersCallback = useCallback(async () => {
    const result = await getUsers();
    setRows(result);
  }, []);
  useEffect(() => {
    getUsersCallback();
  }, []);
  const tableRef = useGridApiRef();

  return (
    <Container sx={{ mt: 2 }}>
      <Stack sx={{ mb: 2 }} direction="row" gap={2}>
        <OpenAddButton />
        <RemoveUserButton tableRef={tableRef} />
      </Stack>
      <DataGrid apiRef={tableRef} columns={columns} rows={rows} checkboxSelection />
    </Container>
  );
};

export default UserList;;

const modalBoxStyle = {
  position: 'absolute' as const,
  top: '50%',
  left: '50%',
  transform: 'translate(-50%, -50%)',
  width: 400,
  bgcolor: 'background.paper',
  boxShadow: 24,
  p: 4,
};
const OpenAddButton = () => {
  const [open, setOpen] = useState(false);
  const { control, register, handleSubmit } = useForm();
  const addUser = useCallback(async (values) => {
    await api.post("/user", values);
    setOpen(false);
    window.location.reload();
  }, []);
  return (
    <>
      <Button variant="contained" onClick={() => setOpen(true)}>添加</Button>
      <Modal open={open}>
        <Box sx={modalBoxStyle}>
          <form onSubmit={handleSubmit(addUser)}>
            <Stack gap={2}>
              <FormTextField control={control} name="Name" required />
              <FormTextField control={control} name="Password" type="password" required />
              <FormRadioGroup control={control} name="Role" options={roleOptions} row valueAsNumber />
              <Stack gap={2} direction="row" justifyContent="center">
                <Button variant="outlined" onClick={() => setOpen(false)}>取消</Button>
                <Button variant="contained" type="submit">添加</Button>
              </Stack>
            </Stack>
          </form>
        </Box>
      </Modal>
    </>
  );
};

interface IRemoveUserButton {
  tableRef: MutableRefObject<GridApiCommunity>;
}
const RemoveUserButton: FC<IRemoveUserButton> = ({ tableRef }) => {

  const removeUser = useCallback(() => {
    const selected = tableRef.current.getSelectedRows();
    selected.forEach(async row => {
      if (row.id && typeof row.id === "string") {
        await api.delete(`/user/${row.id}`);
      }
    });
    window.location.reload();
  }, []);

  return (
    <Button variant="contained" color="error" onClick={removeUser}>刪除</Button>
  );
};;
