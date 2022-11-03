// material-ui
import { alpha } from '@mui/material/styles';
import { Interface } from 'readline';

// ==============================|| DEFAULT THEME - CUSTOM SHADOWS  ||============================== //
export interface CustomShadowsProps {
    button: String;
    text: String;
    z1: String;
}
const CustomShadows = (theme): CustomShadowsProps => ({
    button: `0 2px #0000000b`,
    text: `0 -1px 0 rgb(0 0 0 / 12%)`,
    z1: `0px 2px 8px ${alpha(theme.palette.grey[900], 0.15)}`
});

export default CustomShadows;
