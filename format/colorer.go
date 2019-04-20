package format

var primaryColor = "ffffffff"
var secondaryColor = "ffffffc0"
var warningColor = "ffaa00"
var alarmColor = "ff0000"

// Colorer returns different colors for icon, primary, and
// secondary colors
type Colorer interface {
    IconColor() string
    PrimaryColor() string
    SecondaryColor() string
}

// defaultColorer just returns the default colors {{{
type defaultColorer struct { }

// IconColor returns the default primaryColor
func (d *defaultColorer) IconColor() string {
    return primaryColor
}

// PrimaryColor returns the default primaryColor
func (d *defaultColorer) PrimaryColor() string {
    return primaryColor
}

// SecondaryColor returns the default secondaryColor
func (d *defaultColorer) SecondaryColor() string {
    return secondaryColor
}
// }}}

// dimColorer just returns the default secondaryColor {{{
type dimColorer struct { }

// IconColor returns the default secondaryColor
func (d *dimColorer) IconColor() string {
    return secondaryColor
}

// PrimaryColor returns the default secondaryColor
func (d *dimColorer) PrimaryColor() string {
    return secondaryColor
}

// SecondaryColor returns the default secondaryColor
func (d *dimColorer) SecondaryColor() string {
    return secondaryColor
}
// }}}

// alertColorer returns blinking red {{{
type alertColorer struct { }

// IconColor returns blinking red
func (d *alertColorer) IconColor() string {
    return getAlarmPulseColor()
}

// PrimaryColor returns blinking red
func (d *alertColorer) PrimaryColor() string {
    return getAlarmPulseColor()
}

// SecondaryColor returns blinking red
func (d *alertColorer) SecondaryColor() string {
    return getAlarmPulseColor()
}
// }}}

// warnColorer returns slow blinking orange {{{
type warnColorer struct { }

// IconColor returns slow blinking orange
func (d *warnColorer) IconColor() string {
    return getWarnPulseColor()
}

// PrimaryColor returns slow blinking orange
func (d *warnColorer) PrimaryColor() string {
    return getWarnPulseColor()
}

// SecondaryColor returns slow blinking orange
func (d *warnColorer) SecondaryColor() string {
    return getWarnPulseColor()
}
// }}}

// vim: foldmethod=marker
