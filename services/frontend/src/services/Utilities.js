export function localDate(date) {
    return date.toLocaleString("ru", {
        dateStyle: "medium",
        timeStyle: "long",
    });
}
