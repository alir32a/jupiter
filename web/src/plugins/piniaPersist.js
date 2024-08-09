function persist({ options, store }) {
    const storeID = store.$id;

    const storageState = JSON.parse(window.localStorage.getItem(storeID));

    if (storageState) {
        store.$patch(storageState);
    }

    store.$subscribe((mutation, state) => {
        if (options.persist) {
            window.localStorage.setItem(storeID, JSON.stringify(state));
        }
    });
}

export default persist;