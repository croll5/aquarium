let lastBubbleTime = 0;

/** Code permettant d'afficher de magnifique bulles donnant la sensation de nager parmi les ORCs **/
document.addEventListener('mousemove', function(e) {
    const currentTime = Date.now();
    if (currentTime - lastBubbleTime > 100) { // Délai de 100ms entre chaque bulle
        const bubble = document.createElement('div');
        bubble.classList.add('bubble');
        document.body.appendChild(bubble);

        const size = Math.random() * 20 + 10; // Taille aléatoire entre 10px et 30px
        bubble.style.width = size + 'px';
        bubble.style.height = size + 'px';
        bubble.style.left = e.pageX + 'px';
        bubble.style.top = e.pageY + 'px';

        bubble.addEventListener('animationend', function() {
            bubble.remove();
        });

        lastBubbleTime = currentTime;
    }
});