export interface Relation {
    from: any,
    to?: any,
}

const g = (from, to?) => ({from, to} as Relation);

export const menuItems: Array<Relation> = [
    g("Мероприятия", "/"),
    g("Календарь", "/calendar"),
    g("Профиль", "/profile"),
];
