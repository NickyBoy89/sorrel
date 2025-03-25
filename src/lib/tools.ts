export const toJsDate = (dateString: string): Date => {
    const [year, monthNumber, day] = dateString.split("T")[0].split("-");
    return new Date(Number.parseInt(year), Number.parseInt(monthNumber) - 1, Number.parseInt(day))
}