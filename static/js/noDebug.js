for (i = 0; i < document.scripts.length; i++) {
    if (-1 < document.scripts[i].src.indexOf('chrome-extension:')){
        document.scripts[i].remove();
    }
}