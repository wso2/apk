import React from 'react';
import PropTypes from 'prop-types';
import { useMemo } from 'react';

// material-ui
import { CssBaseline, StyledEngineProvider } from '@mui/material';
import { createTheme, ThemeProvider, ThemeOptions, Shadows } from '@mui/material/styles';

// project import
import Palette from './palette';
import Typography from './typography';
import CustomShadows from './shadows';
import componentsOverride from './overrides';

// ==============================|| DEFAULT THEME - MAIN  ||============================== //
interface ThemeCustomizationProps {
    children: React.ReactNode,
}

export default function ThemeCustomization({ children }: ThemeCustomizationProps) {
    const theme  = Palette('light');

    // eslint-disable-next-line react-hooks/exhaustive-deps
    const themeTypography = Typography(`'Source Code Pro',ui-serif,Georgia,Cambria,Times New Roman,Times,serif`);
    const themeCustomShadows = useMemo(() => CustomShadows(theme), [theme]);

    const themeOptions = useMemo(
        (): ThemeOptions => ({
            breakpoints: {
                values: {
                    xs: 0,
                    sm: 768,
                    md: 1024,
                    lg: 1266,
                    xl: 1536
                }
            },
            direction: 'ltr',
            mixins: {
                toolbar: {
                    minHeight: 60,
                    paddingTop: 8,
                    paddingBottom: 8
                }
            },
            palette: theme.palette,
            typography: themeTypography
        }),
        [theme, themeTypography]
    );

    const themes = createTheme(themeOptions);
    themes.components = componentsOverride(themes);

    return (
        <StyledEngineProvider injectFirst>
            <ThemeProvider theme={themes}>
                <>
                    <CssBaseline />
                    {children}
                </>
            </ThemeProvider>
        </StyledEngineProvider>
    );
}

ThemeCustomization.propTypes = {
    children: PropTypes.node
};
